package drainer

import (
	"context"
	"time"

	"github.com/giantswarm/backoff"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/to"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/policy/v1beta1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"

	"github.com/giantswarm/azure-scheduled-events/pkg/key"
)

func (s *DrainEventHandler) drainNode(ctx context.Context, k8sclient kubernetes.Interface, nodename string) error {
	s.Logger.Debugf(ctx, "Getting node %q for draining", nodename)
	node, err := k8sclient.CoreV1().Nodes().Get(ctx, nodename, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		s.Logger.Debugf(ctx, "Node %q was not found, it was probably already drained and deleted", s.LocalNodeName)
		return nil
	} else if err != nil {
		return microerror.Mask(err)
	}

	s.Logger.Debugf(ctx, "Cordoning node %q", nodename)
	err = cordon(ctx, k8sclient, *node)
	if err != nil {
		return microerror.Mask(err)
	}

	s.Logger.Debugf(ctx, "Evicting pods on node %q", nodename)
	return s.evictPods(ctx, k8sclient, *node)
}

func cordon(ctx context.Context, k8sclient kubernetes.Interface, node corev1.Node) error {
	_, err := k8sclient.CoreV1().Nodes().Patch(ctx, node.GetName(), types.StrategicMergePatchType, []byte(`{"spec":{"unschedulable":true}}`), metav1.PatchOptions{}, "")
	if apierrors.IsNotFound(err) {
		return nil
	}
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (s *DrainEventHandler) evictPods(ctx context.Context, k8sclient kubernetes.Interface, node corev1.Node) error {
	o := func() error {
		var customPods []corev1.Pod
		var kubesystemPods []corev1.Pod
		{

			fieldSelector := fields.SelectorFromSet(fields.Set{
				"spec.nodeName": node.GetName(),
			})
			listOptions := metav1.ListOptions{
				FieldSelector: fieldSelector.String(),
			}
			podList, err := k8sclient.CoreV1().Pods(metav1.NamespaceAll).List(ctx, listOptions)
			if err != nil {
				return microerror.Mask(err)
			}

			for _, pod := range podList.Items {
				if key.IsCriticalPod(pod.Name) {
					// ignore critical pods (api, controller-manager and scheduler)
					// they are static pods so kubelet will recreate them anyway and it can cause other issues
					continue
				}
				if key.IsDaemonSetPod(pod) {
					// ignore daemonSet owned pods
					// daemonSets pod are recreated even on unschedulable node so draining doesn't make sense
					// we are aligning here with community as 'kubectl drain' also ignore them
					continue
				}
				if key.IsEvictedPod(pod) {
					// we don't need to care about already evicted pods
					continue
				}

				if pod.GetNamespace() == "kube-system" {
					kubesystemPods = append(kubesystemPods, pod)
				} else {
					customPods = append(customPods, pod)
				}
			}
		}

		left := len(customPods) + len(kubesystemPods)
		if left == 0 {
			return nil
		}

		if len(customPods) > 0 {
			for _, pod := range customPods {
				err := evict(ctx, k8sclient, pod)
				if IsCannotEvictPod(err) {
					continue
				} else if err != nil {
					return microerror.Mask(err)
				}
			}
		}

		if len(kubesystemPods) > 0 && len(customPods) == 0 {
			for _, pod := range kubesystemPods {
				err := evict(ctx, k8sclient, pod)
				if IsCannotEvictPod(err) {
					continue
				} else if err != nil {
					return microerror.Mask(err)
				}
			}
		}

		return microerror.Maskf(evictionInProgressError, "%d pods still pending eviction, waiting", left)
	}

	err := backoff.RetryNotify(o, backoff.NewConstant(15*time.Minute, 10*time.Second), backoff.NewNotifier(s.Logger, ctx))
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func evict(ctx context.Context, k8sclient kubernetes.Interface, pod corev1.Pod) error {
	eviction := &v1beta1.Eviction{
		ObjectMeta: metav1.ObjectMeta{
			Name:      pod.GetName(),
			Namespace: pod.GetNamespace(),
		},
		DeleteOptions: &metav1.DeleteOptions{
			GracePeriodSeconds: terminationGracePeriod(pod),
		},
	}

	err := k8sclient.PolicyV1beta1().Evictions(eviction.GetNamespace()).Evict(ctx, eviction)
	if IsCannotEvictPod(err) {
		return microerror.Mask(cannotEvictPodError)
	} else if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func terminationGracePeriod(pod corev1.Pod) *int64 {
	var d int64 = 60

	if pod.Spec.TerminationGracePeriodSeconds != nil && *pod.Spec.TerminationGracePeriodSeconds > 0 {
		d = *pod.Spec.TerminationGracePeriodSeconds
	}

	return to.Int64P(d)
}
