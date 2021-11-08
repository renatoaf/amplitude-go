package transport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/renatoaf/amplitude-go/amplitude/data"
)

// HttpTransport wraps upload http based transport strategy to Amplitude V2 api.
type HttpTransport struct {
	url string

	client *http.Client
}

// NewHttpTransport constructs a http transport instance.
func NewHttpTransport(url string, client *http.Client) *HttpTransport {
	return &HttpTransport{
		url:    url,
		client: client,
	}
}

// Send performs a POST request to Amplitude upload API.
func (h HttpTransport) Send(apiKey string, events []*data.Event) error {
	body := &UploadRequestBody{
		ApiKey: apiKey,
		Events: events,
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return err
	}

	response, err := h.client.Post(h.url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	defer response.Body.Close()

	responseBody, _ := ioutil.ReadAll(response.Body)

	switch response.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusAccepted:
		return nil

	case http.StatusBadRequest, http.StatusRequestEntityTooLarge, http.StatusTooManyRequests:
		return fmt.Errorf("request error code %v: %v", response.StatusCode, string(responseBody))
	}

	return fmt.Errorf("unexpected error code %v: %v", response.StatusCode, string(responseBody))
}
