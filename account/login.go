package account

import (
	"encoding/xml"
	"io"

	"github.com/space307/go-lmax-api/model"
	"github.com/space307/go-lmax-api/version"
)

const (
	// CfdLive is a live session
	CfdLive = "CFD_LIVE"
	// CfdDemo is a demo session
	CfdDemo = "CFD_DEMO"
)

const (
	loginURI = "/public/security/login"
)

type (
	// CFDType ...
	CFDType = string
	// LoginRequest ...
	LoginRequest struct {
		Username        string  `xml:"username"`
		Password        string  `xml:"password"`
		CFDType         CFDType `xml:"productType"`
		ProtocolVersion string  `xml:"protocolVersion"`
	}
	// LoginResponse ...
	LoginResponse interface {
		OnSuccess(session model.Session)
		OnFailure(code int, reader io.Reader)
	}
)

// NewLoginRequest creates login request object
func NewLoginRequest(username, password string, cfdType CFDType) *LoginRequest {
	return &LoginRequest{
		Username:        username,
		Password:        password,
		CFDType:         cfdType,
		ProtocolVersion: version.ProtocolVersion,
	}
}

// RequestURI ...
func (lr *LoginRequest) RequestURI() string {
	return loginURI
}

func (lr *LoginRequest) Header() (params map[string]string) {
	params = make(map[string]string)
	return
}

// Write ...
func (lr *LoginRequest) Write(writer io.Writer) error {
	rw := model.NewRequestWrapper(lr)
	doc, err := xml.Marshal(rw)
	if err != nil {
		return err
	}
	_, err = writer.Write(doc)
	return err
}

// AddParam ...
func (lr *LoginRequest) AddParam(_, _ string) {

}
