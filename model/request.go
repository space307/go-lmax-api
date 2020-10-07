package model

import (
	"encoding/xml"
	"io"
)

type (
	// Header ...
	Header interface {
		Get(key string) string
	}
	// Request ...
	Request interface {
		Header() map[string]string

		RequestURI() string
		Write(w io.Writer) error

		AddParam(key, value string)
	}
	// RequestWrapper ...
	RequestWrapper struct {
		XMLName xml.Name `xml:"req"`

		Body interface{} `xml:"body"`
	}
)

// NewRequestWrapper ...
func NewRequestWrapper(r interface{}) *RequestWrapper {
	return &RequestWrapper{Body: r}
}

// NewSubscriptionWrapper ...
func NewSubscriptionWrapper(s interface{}) *RequestWrapper {
	return &RequestWrapper{Body: s}
}
