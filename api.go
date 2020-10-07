package lmax

import (
	"bytes"
	"io"
	"net/http"

	"github.com/space307/go-lmax-api/account"
	"github.com/space307/go-lmax-api/core"
	"github.com/space307/go-lmax-api/model"
)

type (
	// API ...
	API struct {
		httpInvoker *core.HttpInvoker
	}
)

// NewAPI ...
func NewAPI(addr string) *API {
	return &API{
		httpInvoker: core.NewHttpInvoker(addr),
	}
}

func (api *API) Login(msg *account.LoginRequest, callback account.LoginResponse) error {
	args := bytes.NewBuffer(nil)
	err := msg.Write(args)
	if err != nil {
		return err
	}
	return api.httpInvoker.Post(msg.RequestURI(), msg.Header(), args, func(code int, header model.Header, reader io.Reader) error {
		if code != http.StatusOK {
			callback.OnFailure(code, reader)
			return nil
		}

		session, err := core.NewSession(header, reader, api.httpInvoker)
		if err != nil {
			return err
		}
		callback.OnSuccess(session)
		return nil
	})
}
