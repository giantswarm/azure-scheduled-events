package drainer

import (
	"context"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
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
	if event.EventType == "Terminate" && event.ResourceType == "VirtualMachine" {
		s.Logger.LogCtx(ctx, "message", "found Terminate event, start draining the node")

		// Drain the node.
		err := s.drainNode(ctx, s.K8sClient, s.LocalNodeName)
		if IsEvictionInProgress(err) {
			s.Logger.LogCtx(ctx, "level", "warning", "message", "node %q not drained in time.", s.LocalNodeName)
		} else if err != nil {
			return microerror.Mask(err)
		} else {
			s.Logger.LogCtx(ctx, "level", "warning", "message", "node %q drained successfully.", s.LocalNodeName)
		}

		// ACK the event to complete termination.
		err = s.AzureMetadataClient.AckEvent(event.EventId)
		if err != nil {
			return microerror.Mask(err)
		}
		s.Logger.LogCtx(ctx, "message", "acked event")
	}

	return nil
}
