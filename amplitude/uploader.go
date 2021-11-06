package amplitude

import (
	"context"

	"github.com/renatoaf/amplitude-go/amplitude/data"
	"github.com/renatoaf/amplitude-go/amplitude/transport"
)

// UploadCallback function executes deferred code during the upload.
type UploadCallback func(*Uploader)

// UploadBatchDelegate function executes when an upload batch happens (succeeded or not).
type UploadBatchDelegate func(*Uploader, []*data.Event, error)

// Uploader represents an upload worker.
type Uploader struct {
	id        int
	apiKey    string
	queue     chan []*data.Event
	transport transport.Transport
	delegate  UploadBatchDelegate
	callback  UploadCallback
}

func (u *Uploader) Start() (context.Context, context.CancelFunc) {
	if u.queue != nil {
		panic("uploader cannot be started twice")
	}

	u.queue = make(chan []*data.Event, 1)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			case batch, valid := <-u.queue:
				if valid {
					u.upload(batch)
				}
			}
		}
	}()

	return ctx, cancel
}

func (u *Uploader) upload(batch []*data.Event) {
	if u.callback != nil {
		defer u.callback(u)
	}

	err := u.transport.Send(u.apiKey, batch)

	if u.delegate != nil {
		u.delegate(u, batch, err)
	}
}
