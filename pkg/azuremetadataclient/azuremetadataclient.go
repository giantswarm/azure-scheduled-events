package azuremetadataclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/giantswarm/microerror"
)

const (
	infoEndpoint   = "http://169.254.169.254/metadata/instance?api-version=2019-08-01"
	eventsEndpoint = "http://169.254.169.254/metadata/scheduledevents?api-version=2019-08-01"

	ackEventBody = `{
	"StartRequests" : [
		{
			"EventId": %s
		}
	]
}`
)

type Client struct {
	httpClient          *http.Client
	localInstanceVMName string
}

type Config struct {
	// Optional http client to be used for HTTP requests.
	HttpClient *http.Client
}

func New(config Config) (*Client, error) {
	httpClient := config.HttpClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: time.Second * 120}
	}

	metadata, err := getInstanceMetadata(httpClient)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return &Client{
		httpClient:          httpClient,
		localInstanceVMName: metadata.Compute.Name,
	}, nil
}

func (am *Client) AckEvent(eventID string) error {
	req, err := http.NewRequest("POST", eventsEndpoint, bytes.NewBuffer([]byte(fmt.Sprintf(ackEventBody, eventID))))
	if err != nil {
		return microerror.Mask(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Metadata", "True")

	resp, err := am.httpClient.Do(req)
	if err != nil {
		return microerror.Mask(err)
	}
	defer resp.Body.Close()

	return nil
}

func (am *Client) FetchEvents() ([]ScheduledEvent, error) {
	response := ScheduledEventResponse{}
	req, err := http.NewRequest("GET", eventsEndpoint, nil)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	req.Header.Add("Metadata", "true")

	res, err := am.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	var filtered []ScheduledEvent

	for _, event := range response.Events {
		// Check if event is related to the local instance.
		for _, resource := range event.Resources {
			if resource == am.localInstanceVMName {
				filtered = append(filtered, event)
				break
			}
		}
	}

	return filtered, nil
}

func (am *Client) GetInstanceMetadata() (*InstanceResponse, error) {
	return getInstanceMetadata(am.httpClient)
}

func getInstanceMetadata(httpClient *http.Client) (*InstanceResponse, error) {
	req, err := http.NewRequest("GET", infoEndpoint, nil)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	req.Header.Add("Metadata", "true")
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, microerror.Maskf(unexpectedStatusCodeError, "expected HTTP status 200, got %v", resp.StatusCode)
	}

	ret := &InstanceResponse{}

	err = json.NewDecoder(resp.Body).Decode(ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}
