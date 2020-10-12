package heartbeat

import (
	"encoding/xml"
	"io"

	"github.com/space307/go-lmax-api/model"
)

const (
	requestURI = "/secure/read/heartbeat"
)

type (
	// Request ...
	Request struct {
		header map[string]string

		Token string `xml:"token"`
	}
)

// NewRequest ...
func NewRequest() *Request {
	return &Request{header: make(map[string]string)}
}

// Header ...
func (r *Request) Header() map[string]string {
	return r.header
}

// RequestURI ...
func (r *Request) RequestURI() string {
	return requestURI
}

// Write ...
func (r *Request) Write(w io.Writer) error {
	rw := model.NewRequestWrapper(r)
	doc, err := xml.Marshal(rw)
	if err != nil {
		return err
	}
	_, err = w.Write(doc)
	return err
}

// AddParam ...
func (r *Request) AddParam(key, value string) {
	r.header[key] = value
}

// PostHeartbeat ...
func PostHeartbeat(
	session model.Session,
	onSuccess func(),
	onFailure func(code int, reader io.Reader),
	onDisconnected func()) error {
	r := NewRequest()
	return session.HeartbeatRequest(r, onSuccess, onFailure, onDisconnected)
}
