package eventhandler

import (
	"context"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"k8s.io/client-go/kubernetes"

	"github.com/giantswarm/azure-scheduled-events/pkg/azuremetadataclient"
	"github.com/giantswarm/azure-scheduled-events/pkg/drain"
)

type DrainEventHandler struct {
	Drainer   drain.Drainer
	K8sClient kubernetes.Interface
	Logger    micrologger.Logger

	AzureMetadataClient *azuremetadataclient.Client
}

func NewDrainEventHandler(drainer drain.Drainer, logger micrologger.Logger, client *azuremetadataclient.Client, k8sclient kubernetes.Interface) *DrainEventHandler {
	return &DrainEventHandler{
		Drainer:   drainer,
		K8sClient: k8sclient,
		Logger:    logger,

		AzureMetadataClient: client,
	}
}

func (s *DrainEventHandler) HandleEvent(ctx context.Context, event azuremetadataclient.ScheduledEvent) error {
	if event.EventType == "Terminate" && event.ResourceType == "VirtualMachine" {
		s.Logger.LogCtx(ctx, "message", "found Terminate event, start draining the node")
		err := s.Drainer(ctx, s.K8sClient, event.Resources[0])
		if err != nil {
			return microerror.Mask(err)
		}

		// TODO we have to wait for the node to be drained before we ACK the event, or the node will be terminated too early.
		err = s.AzureMetadataClient.AckEvent(event.EventId)
		if err != nil {
			return microerror.Mask(err)
		}
		s.Logger.LogCtx(ctx, "message", "drained node and acked event")
	}

	return nil
}
