package drain

import (
	"context"
	"testing"

	"github.com/giantswarm/to"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/azure-scheduled-events/pkg/unittest"
)

func Test(t *testing.T) {
	ctx := context.Background()

	// Given two pods scheduled in the node being terminated.
	nodeBeingTerminated := "node1"
	node := &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: nodeBeingTerminated,
		},
	}
	pod1 := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-pod1",
			Namespace: "default",
		},
		Spec: corev1.PodSpec{
			NodeName:                      nodeBeingTerminated,
			TerminationGracePeriodSeconds: to.Int64P(1),
		},
	}
	pod2 := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-pod2",
			Namespace: "default",
		},
		Spec: corev1.PodSpec{
			NodeName:                      nodeBeingTerminated,
			TerminationGracePeriodSeconds: to.Int64P(1),
		},
	}
	k8sclients := unittest.FakeK8sClient(node, pod1, pod2)

	err := Drain(ctx, k8sclients.K8sClient(), nodeBeingTerminated)
	if err != nil {
		t.Fatal(err)
	}

	// Then no pods should be running on the node.
	pod, err := k8sclients.K8sClient().CoreV1().Pods(pod1.GetNamespace()).Get(pod1.GetName(), metav1.GetOptions{})
	if err != nil {
		if !apierrors.IsNotFound(err) {
			t.Fatal(err)
		}
	} else {
		t.Fatalf("Pod %#q hasn't been evicted", pod.GetName())
	}

	pod, err = k8sclients.K8sClient().CoreV1().Pods(pod2.GetNamespace()).Get(pod2.GetName(), metav1.GetOptions{})
	if err != nil {
		if !apierrors.IsNotFound(err) {
			t.Fatal(err)
		}
	} else {
		t.Fatalf("Pod %#q hasn't been evicted", pod.GetName())
	}
}
