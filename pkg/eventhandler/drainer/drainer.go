package drainer

import (
	"context"
	"fmt"
	"time"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/giantswarm/azure-scheduled-events/pkg/azuremetadataclient"
)

type DrainEventHandler struct {
	K8sClient kubernetes.Interface
	Logger    micrologger.Logger

	AzureMetadataClient *azuremetadataclient.Client
	LocalNodeName       string
}

func NewDrainEventHandler(logger micrologger.Logger, client *azuremetadataclient.Client, k8sclient kubernetes.Interface, localNodeName string) *DrainEventHandler {
	return &DrainEventHandler{
		K8sClient: k8sclient,
		Logger:    logger,

		AzureMetadataClient: client,
		LocalNodeName:       localNodeName,
	}
}

func (s *DrainEventHandler) HandleEvent(ctx context.Context, event azuremetadataclient.ScheduledEvent) error {
	s.Logger.Debugf(ctx, "Received event: %v", event)

	var timeout time.Duration
	{
		t, err := time.Parse(time.RFC1123, event.NotBefore)
		if err != nil {
			timeout = 15 * time.Minute
		} else {
			timeout = time.Until(t)
		}
	}

	switch event.EventType {
	case "Terminate":
	case "Preempt":
	case "Reboot":
	case "Redeploy":
	default:
		s.Logger.Debugf(ctx, "Warning: unhandled event type: %s (resource type %s)", event.EventType, event.ResourceType)
		return nil
	}

	// If timeout is negative we have no time to drain the node, so we skip the step.
	if timeout > 0 {
		s.Logger.Debugf(ctx, "got event %s, start draining the node (timeout %.0f seconds)", event.EventType, timeout.Seconds())

		// Drain the node.
		err := s.drainNode(ctx, s.K8sClient, s.LocalNodeName, timeout)
		if IsEvictionInProgress(err) {
			s.Logger.LogCtx(ctx, "level", "warning", "message", fmt.Sprintf("node %q not drained in time.", s.LocalNodeName))
		} else if err != nil {
			return microerror.Mask(err)
		} else {
			s.Logger.Debugf(ctx, fmt.Sprintf("node %q drained successfully.", s.LocalNodeName))
		}
	}

	// Delete node from k8s.
	s.Logger.Debugf(ctx, "Deleting Node %q from k8s API", s.LocalNodeName)
	err := s.K8sClient.CoreV1().Nodes().Delete(ctx, s.LocalNodeName, v1.DeleteOptions{})
	if apierrors.IsNotFound(err) {
		s.Logger.Debugf(ctx, "Node %q was not found, it was probably already deleted", s.LocalNodeName)
	} else if err != nil {
		s.Logger.LogCtx(ctx, "level", "warning", "message", fmt.Sprintf("Error deleting node %q from Kubernetes API: %s.", s.LocalNodeName, err))
	}

	// ACK the event.
	err = s.AzureMetadataClient.AckEvent(event.EventId)
	if err != nil {
		return microerror.Mask(err)
	}
	s.Logger.LogCtx(ctx, "message", fmt.Sprintf("acked event %q", event.EventId))

	return nil
}
