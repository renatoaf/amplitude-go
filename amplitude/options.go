package amplitude

import (
	"context"
	"time"

	"github.com/renatoaf/amplitude-go/amplitude/transport"
)

// Options defines API Client optional parameters.
type Options struct {
	// Base context
	Context context.Context

	// BatchSize is the amount of events sent in a single upload batch; defaults to DefaultBatchSize.
	BatchSize int

	// MaxCachedEvents is the limit of events we can queue before start dropping logged events; defaults to DefaultMaxCachedEvents.
	MaxCachedEvents int

	// MaxParallelUploads is the number of goroutines we spawn to handle uploads in batch; defaults to DefaultMaxParallelUploads.
	MaxParallelUploads int

	// MaxUploadAttempts is the number of upload attempts we execute before dropping a batch (MaxUploadAttempts - 1 retries); defaults to DefaultMaxUploadAttempts.
	MaxUploadAttempts int

	// OptOut should be set to true if you want a dry-run experience; defaults to false.
	OptOut bool

	// Transport allows you to configure your own transport layer; defaults to transport.HttpTransport.
	Transport transport.Transport

	// ServerUrl allows you to configure your own API server url; defaults to DefaultServerUrl.
	ServerUrl string

	// RequestTimeout is the timeout for each individual upload attempt; defaults to DefaultRequestTimeout.
	RequestTimeout time.Duration

	// UploadInterval is the interval between individual uploads; defaults to DefaultUploadInterval.
	UploadInterval time.Duration

	// UploadDelegate allows you to hook your own code when an upload happens (for example to send metrics); defaults to nil.
	UploadDelegate UploadBatchDelegate // defaults to nil
}

func (o *Options) BaseContext() context.Context {
	if o.Context != nil {
		return o.Context
	}

	return context.Background()
}
