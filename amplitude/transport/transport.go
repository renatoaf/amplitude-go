package transport

import (
	"github.com/renatoaf/amplitude-go/amplitude/data"
)

// Transport defines a send strategy.
type Transport interface {
	// Send uploads the events to a specific api key.
	Send(apiKey string, events []*data.Event) error
}
