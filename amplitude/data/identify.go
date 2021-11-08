package data

import "encoding/json"

// Special event type for identify events.
const identifySpecialEventType = "$identify"

// TODO support "identify properties"
// Identify defines the properties of an $identify event on HTTP v2 api: https://developers.amplitude.com/docs/http-api-v2#properties-1.
type Identify struct {
	UserID         string
	DeviceID       string
	UserProperties map[string]interface{}
}

func (i Identify) AsEvent() *Event {
	return &Event{
		UserID:         i.UserID,
		DeviceID:       i.DeviceID,
		UserProperties: i.UserProperties,
		EventType:      identifySpecialEventType,
	}
}

func (i Identify) String() string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}
