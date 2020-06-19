package scheduledevents

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"k8s.io/client-go/kubernetes"

	"github.com/giantswarm/azure-scheduled-events/pkg/drain"
)

type ScheduledEvent struct {
	EventId      string   `json:"eventid"`
	EventType    string   `json:"eventtype"`
	ResourceType string   `json:"resourcetype"`
	Resources    []string `json:"resources"`
	EventStatus  string   `json:"eventstatus"`
	NotBefore    string   `json:"notbefore"`
}

type MetadataResponse struct {
	DocumentIncarnation string           `json:"incarnation"`
	Events              []ScheduledEvent `json:"events"`
}

const (
	DefaultMetadataEndpoint = "http://169.254.169.254/metadata/scheduledevents?api-version=2019-01-01"
	ackEventBody            = `{
	"StartRequests" : [
		{
			"EventId": %s
		}
	]
}`
)

type ScheduledEvents struct {
	Drainer drain.Drainer
	Logger  micrologger.Logger
}

func NewScheduledEvents(drainer drain.Drainer, logger micrologger.Logger) ScheduledEvents {
	return ScheduledEvents{
		Drainer: drainer,
		Logger:  logger,
	}
}

func fetchEvents(metadataURL string) (MetadataResponse, error) {
	response := MetadataResponse{}
	res, err := http.Get(metadataURL)
	if err != nil {
		return response, err
	}

	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return response, err
	}

	return response, nil
}

func (s *ScheduledEvents) GetEvents(ctx context.Context, k8sclient kubernetes.Interface, metadataURL string) error {
	s.Logger.LogCtx(ctx, "message", "fetching events from metadata endpoint")
	response, err := fetchEvents(metadataURL)
	if err != nil {
		return microerror.Mask(err)
	}

	s.Logger.LogCtx(ctx, "message", "looping through received events")
	for _, event := range response.Events {
		if event.EventType == "Terminate" && event.ResourceType == "VirtualMachine" {
			s.Logger.LogCtx(ctx, "message", "found terminated event, let's drain")
			err = s.Drainer(ctx, k8sclient, event.Resources[0])
			if err != nil {
				return microerror.Mask(err)
			}

			err = ackEvent(event.EventId, metadataURL)
			if err != nil {
				return microerror.Mask(err)
			}
			s.Logger.LogCtx(ctx, "message", "drained node and acked event")
		}
	}

	s.Logger.LogCtx(ctx, "message", "finished looping through events")

	return nil
}

func ackEvent(eventId, metadataURL string) error {
	req, err := http.NewRequest("POST", metadataURL, bytes.NewBuffer([]byte(fmt.Sprintf(ackEventBody, eventId))))
	if err != nil {
		return microerror.Mask(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Metadata", "True")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	return nil
}
