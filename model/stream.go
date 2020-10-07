package model

type (
	// Stream ...
	Stream interface {
		Serve() error
		Stop()
	}

	// StreamHandler ...
	StreamHandler interface {
		Handle(data []byte)
	}
)
