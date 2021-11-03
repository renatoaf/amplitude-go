package amplitude

import "time"

const (
	DefaultServerUrl = "https://api2.amplitude.com/2/httpapi"

	// 256 events batch size for individual uploads
	DefaultBatchSize = 256

	// Based on nodejs client (but not estimated for this SDK): 2kb is a safe estimate for a medium size event object. This keeps the SDK's memory footprint roughly under 32 MB.
	DefaultMaxCachedEvents = 16000

	// 12 goroutines to execute parallel uploads.
	DefaultMaxParallelUploads = 12

	// 5 retries with backoff.
	DefaultMaxUploadAttempts = 6

	// 10 ms pace of upload
	DefaultUploadInterval = 10 * time.Millisecond

	// Default of 10s was took from NodeJS sdk
	DefaultRequestTimeout = 10 * time.Second
)
