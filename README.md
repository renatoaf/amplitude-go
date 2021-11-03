# Amplitude Golang SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/renatoaf/amplitude-go.svg)](https://pkg.go.dev/github.com/renatoaf/amplitude-go)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/renatoaf/amplitude-go)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/renatoaf/amplitude-go)

Amplitude unofficial client for Go, inspired in their official [SDK for Node](https://github.com/amplitude/Amplitude-Node).

For reference, visit HTTP API v2 [documentation](https://developers.amplitude.com/docs/http-api-v2).

## Installation

	$ go get github.com/renatoaf/amplitude-go

## Usage

```go
// startup
client := amplitude.NewDefaultClient("<your-api-key")
client.Start()

// logging events
client.LogEvent(&data.Event{
    UserID: "datamonster@gmail.com",
    EventType: "test-event",
    EventProperties: map[string]interface{}{
        "source": "notification",
    },
    UserProperties: map[string]interface{}{
        "age": 25,
        "gender": "female",
    },
})

// gracefully shutdown, waiting pending events to be sent
client.Shutdown()
```

The `Event` ([doc](https://pkg.go.dev/github.com/renatoaf/amplitude-go/amplitude/data#Event)) structure is based on API V2 [request properties](https://developers.amplitude.com/docs/http-api-v2#properties-1).

Events will not be sent synchronously, the client keeps a goroutine responsible for batching and issuing uploads of events.

The client will upload events:
  - after the upload interval (every 10ms by default).
  - as soon as we accumulate enough events to batch (256 events by default).
  - when `Flush` is explicitly invoked.
  - during shutdown process.

`LogEvent` should never block, it will return an error in case the event was not queued (which means the event will be dropped without even being sent). This should not happen unless the uploads are not getting through.

Check advanced parameters to learn how to tweak the default behaviour.

### Advanced parameters

The default client behaviour can be configured through a set of custom `Options` ([doc](https://pkg.go.dev/github.com/renatoaf/amplitude-go/amplitude#Options)).

```go
client := amplitude.NewClient("<your-api-key", amplitude.Options{ ... })
```

### Examples

1. If you want to configure your client to issue uploads every second:

```go
client := amplitude.NewClient("<your-api-key", amplitude.Options{
    UploadInterval: time.Second,
})
```

2. If you want to disable retries:

```go
client := amplitude.NewClient("<your-api-key", amplitude.Options{
    MaxUploadAttempts: 1,
})
```

3. If you want to hook your own Datadog metrics for amplitude events:

```go
client := amplitude.NewClient("<your-api-key", amplitude.Options{
    UploadDelegate: func(_ *amplitude.Uploader, events []*data.Event, err error) {
        count := len(events)
        if err != nil {
            statsd.Incr("amplitude.events", []string{"status:failure"}, count)
        } else {
            statsd.Incr("amplitude.events", []string{"status:success"}, count)
        }
    },
})
```