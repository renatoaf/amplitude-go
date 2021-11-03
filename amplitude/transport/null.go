package transport

import (
	"github.com/renatoaf/amplitude-go/amplitude/data"
)

// NullTransport represents a null transport strategy.
type NullTransport struct {
}

// Send in a null transport strategy does nothing.
func (n NullTransport) Send(apiKey string, events []*data.Event) error {
	return nil
}
