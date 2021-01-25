package drainer

//import (
//	"context"
//	"encoding/json"
//	"fmt"
//	"net/http"
//	"net/http/httptest"
//	"testing"
//
//	"github.com/giantswarm/micrologger"
//	"github.com/giantswarm/to"
//	corev1 "k8s.io/api/core/v1"
//	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
//
//	"github.com/giantswarm/azure-scheduled-events/pkg/azuremetadataclient"
//	"github.com/giantswarm/azure-scheduled-events/pkg/unittest"
//)
//
//func Test(t *testing.T) {
//	ctx := context.Background()
//
//	// Given two pods scheduled in the node being terminated.
//	nodeBeingTerminated := "node1"
//	node := &corev1.Node{
//		ObjectMeta: metav1.ObjectMeta{
//			Name: nodeBeingTerminated,
//		},
//	}
//	pod1 := &corev1.Pod{
//		ObjectMeta: metav1.ObjectMeta{
//			Name:      "my-pod1",
//			Namespace: "default",
//		},
//		Spec: corev1.PodSpec{
//			NodeName:                      nodeBeingTerminated,
//			TerminationGracePeriodSeconds: to.Int64P(1),
//		},
//	}
//	pod2 := &corev1.Pod{
//		ObjectMeta: metav1.ObjectMeta{
//			Name:      "my-pod2",
//			Namespace: "default",
//		},
//		Spec: corev1.PodSpec{
//			NodeName:                      nodeBeingTerminated,
//			TerminationGracePeriodSeconds: to.Int64P(1),
//		},
//	}
//	k8sclients := unittest.FakeK8sClient(node, pod1, pod2)
//
//	logger,err := micrologger.New(micrologger.Config{})
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		resp := azuremetadataclient.InstanceResponse{Compute: azuremetadataclient.Compute{
//			Name: nodeBeingTerminated,
//		}}
//
//		body, err := json.Marshal(resp)
//		if err != nil {
//			t.Fatal(err)
//		}
//
//		fmt.Fprintln(w, string(body))
//	}))
//	defer ts.Close()
//
//	mc,err := azuremetadataclient.New(azuremetadataclient.Config{ HttpClient: ts.Client() })
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	s := DrainEventHandler{
//		K8sClient:           k8sclients.K8sClient(),
//		Logger:              logger,
//		AzureMetadataClient: mc,
//		LocalNodeName:       nodeBeingTerminated,
//	}
//
//	err = s.drainNode(ctx, k8sclients.K8sClient(), nodeBeingTerminated)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	// Then no pods should be running on the node.
//	//pod, err := k8sclients.K8sClient().CoreV1().Pods(pod1.GetNamespace()).Get(pod1.GetName(), metav1.GetOptions{})
//	//if err != nil {
//	//	if !apierrors.IsNotFound(err) {
//	//		t.Fatal(err)
//	//	}
//	//} else {
//	//	t.Fatalf("Pod %#q hasn't been evicted", pod.GetName())
//	//}
//
//	//pod, err = k8sclients.K8sClient().CoreV1().Pods(pod2.GetNamespace()).Get(pod2.GetName(), metav1.GetOptions{})
//	//if err != nil {
//	//	if !apierrors.IsNotFound(err) {
//	//		t.Fatal(err)
//	//	}
//	//} else {
//	//	t.Fatalf("Pod %#q hasn't been evicted", pod.GetName())
//	//}
//}
