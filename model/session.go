package model

import (
	"io"

	"github.com/space307/go-lmax-api/events"
)

type (
	// RestAPI ...
	RestAPI interface {
		Post(request Request, success func(r io.Reader), failure func(code int, r io.Reader)) error
		Get(request Request, success func(r io.Reader), failure func(code int, r io.Reader)) error
	}

	// SessionCredentials ...
	SessionCredentials interface {
		ID() string
		UserAgent() string
	}

	// EventsManager ...
	EventsManager interface {
		AddEventListener(t events.Type, o events.Observer)
		RemoveEventListener(t events.Type, o events.Observer)
	}

	// Session ...
	Session interface {
		Server
		RestAPI
		SessionCredentials
		EventsManager

		Logout(request Request, success func(r io.Reader), failure func(code int, r io.Reader)) error
	}
)
