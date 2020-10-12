package model

import (
	"encoding/xml"
	"io"
)

const (
	subscriptionURI = "/secure/subscribe"
)

const (
	// AccountStateSubscription ...
	AccountStateSubscription   = "account"
	PositionsStateSubscription = "position"
	OrdersStateSubscription    = "order"
	HeartbeatSubscription      = "account"
)

type (
	// Subscription ...
	Subscription struct {
		XMLName xml.Name `xml:"body"`

		Body struct {
			XMLName xml.Name `xml:"subscription"`

			Type string `xml:"type"`
		}

		header map[string]string
	}
)

// NewStateRequest ...
func NewSubscription() *Subscription {
	return &Subscription{header: make(map[string]string)}
}

// PositionsSubscriptionRequest ...
func (s *Subscription) Write(w io.Writer) error {
	rw := NewSubscriptionWrapper(s)
	doc, err := xml.Marshal(rw)
	if err != nil {
		return err
	}
	_, err = w.Write(doc)
	return err
}

// RequestURI ...
func (s *Subscription) RequestURI() string {
	return subscriptionURI
}

// Header ...
func (s *Subscription) Header() map[string]string {
	return s.header
}

// Header ...
func (s *Subscription) AddParam(key, value string) {
	s.header[key] = value
}

// SetType ...
func (s *Subscription) SetType(t string) {
	s.Body.Type = t
}
