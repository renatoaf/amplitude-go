package data

import "encoding/json"

// Event defines the event properties on HTTP v2 api: https://developers.amplitude.com/docs/http-api-v2#properties-1.
type Event struct {
	UserID             string                 `json:"user_id,omitempty"`
	DeviceID           string                 `json:"device_id,omitempty"`
	EventType          string                 `json:"event_type,omitempty"`
	Time               int64                  `json:"time,omitempty"`
	EventProperties    map[string]interface{} `json:"event_properties,omitempty"`
	UserProperties     map[string]interface{} `json:"user_properties,omitempty"`
	AppVersion         string                 `json:"app_version,omitempty"`
	Platform           string                 `json:"platform,omitempty"`
	OSName             string                 `json:"os_name,omitempty"`
	OSVersion          string                 `json:"os_version,omitempty"`
	DeviceBrand        string                 `json:"device_brand,omitempty"`
	DeviceManufacturer string                 `json:"device_manufacturer,omitempty"`
	DeviceModel        string                 `json:"device_model,omitempty"`
	Carrier            string                 `json:"carrier,omitempty"`
	Country            string                 `json:"country,omitempty"`
	Region             string                 `json:"region,omitempty"`
	City               string                 `json:"city,omitempty"`
	DMA                string                 `json:"dma,omitempty"`
	Language           string                 `json:"language,omitempty"`
	Price              float64                `json:"price,omitempty"`
	Quantity           int32                  `json:"quantity,omitempty"`
	Revenue            float64                `json:"revenue,omitempty"`
	ProductID          string                 `json:"productId,omitempty"`
	RevenueType        string                 `json:"revenueType,omitempty"`
	Latitude           float64                `json:"location_lat,omitempty"`
	Longitude          float64                `json:"location_lng,omitempty"`
	IP                 string                 `json:"ip,omitempty"`
	IDFA               string                 `json:"idfa,omitempty"`
	IDFV               string                 `json:"idfv,omitempty"`
	ADID               string                 `json:"adid,omitempty"`
	AndroidID          string                 `json:"android_id,omitempty"`
	EventID            int32                  `json:"event_id,omitempty"`
	SessionID          int64                  `json:"session_id,omitempty"`
	InsertID           string                 `json:"insert_id,omitempty"`
}

func (e Event) String() string {
	s, _ := json.MarshalIndent(e, "", "\t")
	return string(s)
}
