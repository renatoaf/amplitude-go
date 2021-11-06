package transport

import (
	"encoding/json"

	"github.com/renatoaf/amplitude-go/amplitude/data"
)

// UploadRequestBody wraps the upload request body on http v2 api: https://developers.amplitude.com/docs/http-api-v2#schemauploadrequestbody
type UploadRequestBody struct {
	ApiKey string        `json:"api_key,omitempty"`
	Events []*data.Event `json:"events,omitempty"`
}

func (u UploadRequestBody) String() string {
	s, _ := json.MarshalIndent(u, "", "\t")
	return string(s)
}
