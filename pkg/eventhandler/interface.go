package eventhandler

import (
	"context"

	"github.com/giantswarm/azure-scheduled-events/pkg/azuremetadataclient"
)

type EventHandler interface {
	HandleEvent(ctx context.Context, event azuremetadataclient.ScheduledEvent) error
}
