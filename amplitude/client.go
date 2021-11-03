package amplitude

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/renatoaf/amplitude-go/amplitude/data"
	"github.com/renatoaf/amplitude-go/amplitude/transport"
	"net/http"
	"sync"
	"time"
)

// ClientState defines the client FSM.
type ClientState uint32

const (
	Idle ClientState = iota
	Running
	Closing
	Closed
)

// Client wraps Amplitude V2 api client.
type Client struct {
	state ClientState

	ctx    context.Context
	cancel context.CancelFunc

	writes  *sync.WaitGroup
	uploads *sync.WaitGroup

	workers  chan *Uploader
	shutdown chan interface{}
	flush    chan chan interface{}
	closing  chan interface{}
	buffer   chan *data.Event
	timer    *time.Timer

	apiKey  string
	options Options
}

// NewClient creates a client with a set of options.
func NewClient(apiKey string, options Options) *Client {
	if apiKey == "" {
		panic("api key must be set")
	}

	sanitizeOptions(&options)

	return &Client{
		state:   Idle,
		writes:  &sync.WaitGroup{},
		uploads: &sync.WaitGroup{},
		apiKey:  apiKey,
		options: options,
	}
}

// NewDefaultClient creates a client with default options.
func NewDefaultClient(apiKey string) *Client {
	return NewClient(apiKey, Options{})
}

// State returns client current state.
func (c *Client) State() ClientState {
	return c.state
}

// Start starts the client event transport routine.
func (c *Client) Start() error {
	if c.state != Idle {
		return fmt.Errorf("client already started")
	}

	c.state = Running
	c.ctx, c.cancel = context.WithCancel(context.Background())
	c.flush = make(chan chan interface{}, 1)
	c.shutdown = make(chan interface{}, 1)
	c.closing = make(chan interface{}, 1)
	c.buffer = make(chan *data.Event, c.options.MaxCachedEvents)
	c.workers = make(chan *Uploader, c.options.MaxParallelUploads)
	c.timer = time.NewTimer(c.options.UploadInterval)

	for i := 0; i < cap(c.workers); i++ {
		uploader := &Uploader{
			id:        i,
			apiKey:    c.apiKey,
			transport: buildTransport(&c.options),
			callback: func(uploader *Uploader) {
				c.workers <- uploader
				c.uploads.Done()
			},
			delegate: c.options.UploadDelegate,
		}
		uploader.Start()

		c.workers <- uploader
	}

	go c.run(c.options.BatchSize)
	return nil
}

// Shutdown stops the client and flushes remaining events.
func (c *Client) Shutdown() error {
	if c.state != Running {
		return fmt.Errorf("client not running")
	}

	c.state = Closing
	c.cancel()

	<-c.shutdown
	return nil
}

// LogEvent enqueues an event to be uploaded in background.
func (c *Client) LogEvent(event *data.Event) error {
	if c.options.OptOut {
		return nil
	}

	if c.state != Running {
		return fmt.Errorf("client not running")
	}

	c.writes.Add(1)
	defer c.writes.Done()

	select {
	case <-c.closing:
		return fmt.Errorf("client is finishing")
	default:
		// ignore
	}

	select {
	case <-c.closing:
		return fmt.Errorf("client is finishing")
	case c.buffer <- event:
		return nil
	default:
		return fmt.Errorf("client queue is full")
	}
}

// Flush uploads all queued events immediately.
func (c *Client) Flush() error {
	if c.options.OptOut {
		return nil
	}

	if c.state != Running {
		return fmt.Errorf("client not running")
	}

	request := make(chan interface{})
	defer close(request)

	c.flush <- request

	<-request // wait for it to complete
	return nil
}

// internal functions.
func (c *Client) run(batchSize int) {
	buffer := make([]*data.Event, 0, batchSize)

	upload := func() {
		if len(buffer) > 0 {
			c.uploads.Add(1)

			uploader := <-c.workers // get a free uploader from the pool or wait.
			uploader.queue <- buffer

			buffer = make([]*data.Event, 0, batchSize)
		}

		c.timer.Reset(c.options.UploadInterval)
	}

	for c.state != Closed {
		select {
		case <-c.ctx.Done():
			c.state = Closed
			break

		case <-c.timer.C:
			upload()

		case flushRequest, valid := <-c.flush:
			if valid {
				upload()

				flushRequest <- true
			}

		case event, valid := <-c.buffer:
			if valid {
				buffer = append(buffer, event)
			}

			if len(buffer) >= batchSize {
				upload()
			}
		}
	}

	// https://www.leolara.me/blog/closing_a_go_channel_written_by_several_goroutines/
	close(c.closing)

	c.writes.Wait()

	close(c.buffer)

	for event := range c.buffer { // consume everything once the channel is closed.
		buffer = append(buffer, event)
	}

	upload()

	c.uploads.Wait()

	c.shutdown <- true
}

func sanitizeOptions(o *Options) {
	if o.ServerUrl == "" {
		o.ServerUrl = DefaultServerUrl
	}

	if o.BatchSize == 0 {
		o.BatchSize = DefaultBatchSize
	}

	if o.MaxCachedEvents == 0 {
		o.MaxCachedEvents = DefaultMaxCachedEvents
	}

	if o.RequestTimeout == 0 {
		o.RequestTimeout = DefaultRequestTimeout
	}

	if o.UploadInterval == 0 {
		o.UploadInterval = DefaultUploadInterval
	}

	if o.MaxParallelUploads == 0 {
		o.MaxParallelUploads = DefaultMaxParallelUploads
	}

	if o.MaxUploadAttempts == 0 {
		o.MaxUploadAttempts = DefaultMaxUploadAttempts
	}
}

func buildTransport(options *Options) transport.Transport {
	uploadTransport := options.Transport
	if uploadTransport == nil {
		uploadTransport = buildDefaultTransport(options)
	}

	return uploadTransport
}

func buildDefaultTransport(options *Options) transport.Transport {
	client := &http.Client{}

	if options.MaxUploadAttempts > 1 {
		retryableClient := retryablehttp.NewClient()
		retryableClient.RetryMax = options.MaxUploadAttempts - 1

		client = retryableClient.StandardClient()
	}

	client.Timeout = options.RequestTimeout

	return transport.NewHttpTransport(options.ServerUrl, client)
}
