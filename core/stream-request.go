package core

import (
	"io"
)

const (
	requestURI = "/push/stream"
)

type (
	// StreamRequest ...
	StreamRequest struct {
		header map[string]string
	}
)

// NewStreamRequest ...
func NewStreamRequest() *StreamRequest {
	return &StreamRequest{header: make(map[string]string)}
}

// StateRequest ...
func (sr *StreamRequest) Write(w io.Writer) error {
	return nil
}

func (sr *StreamRequest) Header() map[string]string {
	return sr.header
}

func (sr *StreamRequest) AddParam(key, value string) {
	sr.header[key] = value
}

// RequestURI ...
func (sr *StreamRequest) RequestURI() string {
	return requestURI
}
