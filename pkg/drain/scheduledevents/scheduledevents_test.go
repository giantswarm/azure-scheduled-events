package scheduledevents

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"k8s.io/client-go/kubernetes"

	"github.com/giantswarm/azure-scheduled-events/pkg/unittest"
)

const (
	jsonResponse = `{
    "DocumentIncarnation": "abc123",
    "Events": [
        {
            "EventId": "123",
            "EventType": "Terminate",
            "ResourceType": "VirtualMachine",
            "Resources": ["nodename"],
            "EventStatus": "Scheduled",
            "NotBefore": "20170702"
        }
    ]
}`
)

type MockDrainer struct {
	nodename string
}

var mockdrainer = MockDrainer{}

func mockDrainer(ctx context.Context, k8sclient kubernetes.Interface, nodename string) error {
	mockdrainer.nodename = nodename
	return nil
}

func Test(t *testing.T) {
	ctx := context.Background()
	k8sclients := unittest.FakeK8sClient()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintln(w, jsonResponse)
		if err != nil {
			t.Fatal(err)
		}
	}))
	defer ts.Close()

	scheduledEvents := NewScheduledEvents(mockDrainer)
	err := scheduledEvents.GetEvents(ctx, k8sclients.K8sClient(), ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	if mockdrainer.nodename != "nodename" {
		t.Fatal("Drain wasn't called or called with wrong nodename")
	}
}
