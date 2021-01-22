package azuremetadata

type (
	ScheduledEventResponse struct {
		DocumentIncarnation int              `json:"DocumentIncarnation"`
		Events              []ScheduledEvent `json:"Events"`
	}

	ScheduledEvent struct {
		EventId      string   `json:"EventId"`
		EventType    string   `json:"EventType"`
		ResourceType string   `json:"ResourceType"`
		Resources    []string `json:"Resources"`
		EventStatus  string   `json:"EventStatus"`
		NotBefore    string   `json:"NotBefore"`
	}

	InstanceResponse struct {
		Compute struct {
			Location             string `json:"location"`
			Name                 string `json:"name"`
			Offer                string `json:"offer"`
			OsType               string `json:"osType"`
			PlacementGroupID     string `json:"placementGroupId"`
			PlatformFaultDomain  string `json:"platformFaultDomain"`
			PlatformUpdateDomain string `json:"platformUpdateDomain"`
			Publisher            string `json:"publisher"`
			ResourceGroupName    string `json:"resourceGroupName"`
			Sku                  string `json:"sku"`
			SubscriptionID       string `json:"subscriptionId"`
			Tags                 string `json:"tags"`
			Version              string `json:"version"`
			VMID                 string `json:"vmId"`
			VMSize               string `json:"vmSize"`
		} `json:"compute"`
	}
)
