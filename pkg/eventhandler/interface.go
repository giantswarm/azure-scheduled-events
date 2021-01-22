package eventhandler

import (
	"context"

	"github.com/giantswarm/azure-scheduled-events/pkg/azuremetadata"
)

type EventHandler interface {
	HandleEvent(ctx context.Context, event azuremetadata.ScheduledEvent) error
}
