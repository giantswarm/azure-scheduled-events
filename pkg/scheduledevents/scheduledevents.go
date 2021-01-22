package scheduledevents

import (
	"context"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"k8s.io/client-go/kubernetes"

	"github.com/giantswarm/azure-scheduled-events/pkg/azuremetadata"
	"github.com/giantswarm/azure-scheduled-events/pkg/drain"
)

type ScheduledEvents struct {
	Drainer drain.Drainer
	Logger  micrologger.Logger

	AzureMetadataClient *azuremetadata.Client
}

func NewScheduledEvents(drainer drain.Drainer, logger micrologger.Logger, client *azuremetadata.Client) ScheduledEvents {
	return ScheduledEvents{
		Drainer: drainer,
		Logger:  logger,

		AzureMetadataClient: client,
	}
}

func (s *ScheduledEvents) GetEvents(ctx context.Context, k8sclient kubernetes.Interface) error {
	s.Logger.LogCtx(ctx, "message", "fetching events from metadata endpoint")
	events, err := s.AzureMetadataClient.FetchEvents()
	if err != nil {
		return microerror.Mask(err)
	}

	s.Logger.LogCtx(ctx, "message", "looping through received events")
	for _, event := range events {
		if event.EventType == "Terminate" && event.ResourceType == "VirtualMachine" {
			s.Logger.LogCtx(ctx, "message", "found terminated event, let's drain")
			err = s.Drainer(ctx, k8sclient, event.Resources[0])
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
	}

	s.Logger.LogCtx(ctx, "message", "finished looping through events")

	return nil
}
