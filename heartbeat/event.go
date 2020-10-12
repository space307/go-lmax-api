package heartbeat

import (
	"io"

	"github.com/space307/go-lmax-api/events"
	"github.com/space307/go-lmax-api/model"
)

type (
	// Heartbeat ...
	Event struct {
		AccountID int64  `xml:"accountId"`
		Token     string `xml:"token"`
	}
)

// Type ...
func (_ *Event) Type() events.Type {
	return events.Heartbeat
}

// SubscribeState ...
func SubscribeState(session model.Session, onSuccess func(reader io.Reader), onFailure func(code int, reader io.Reader)) error {
	subscription := model.NewSubscription()
	subscription.SetType(model.HeartbeatSubscription)
	return session.Post(subscription, onSuccess, onFailure)
}
