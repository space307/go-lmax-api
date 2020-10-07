package model

type (
	// Server ...
	Server interface {
		Serve() error
		Stop() error
		Wait()
	}
)
