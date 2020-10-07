package account

import (
	"encoding/xml"
	"io"

	"github.com/space307/go-lmax-api/model"
)

const (
	requestURI = "/secure/account/requestAccountState"
)

type (
	// StateRequest ...
	StateRequest struct {
		header map[string]string
	}
)

// NewStateRequest ...
func NewStateRequest() *StateRequest {
	return &StateRequest{header: make(map[string]string)}
}

// StateRequest ...
func (sr *StateRequest) Write(w io.Writer) error {
	rw := model.NewRequestWrapper(sr)
	doc, err := xml.Marshal(rw)
	if err != nil {
		return err
	}
	_, err = w.Write(doc)
	return err
}

func (sr *StateRequest) Header() (result map[string]string) {
	return sr.header
}

// RequestURI ...
func (sr *StateRequest) RequestURI() string {
	return requestURI
}

func (sr *StateRequest) AddParam(key, value string) {
	sr.header[key] = value
}

// GetState ...
func GetState(s model.Session, onSuccess func(reader io.Reader), onFailure func(code int, reader io.Reader)) error {
	sr := NewStateRequest()
	return s.Get(sr, onSuccess, onFailure)
}
