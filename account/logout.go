package account

import (
	"encoding/xml"
	"io"

	"github.com/space307/go-lmax-api/model"
)

const (
	logoutURI = "/public/security/logout"
)

type (
	LogoutRequest struct {
	}
)

// NewLogoutRequest creates login request object
func NewLogoutRequest() *LogoutRequest {
	return &LogoutRequest{}
}

// RequestURI ...
func (lr *LogoutRequest) RequestURI() string {
	return logoutURI
}

func (lr *LogoutRequest) Header() (params map[string]string) {
	params = make(map[string]string)
	return
}

// Write ...
func (lr *LogoutRequest) Write(writer io.Writer) error {
	rw := model.NewRequestWrapper(lr)
	doc, err := xml.Marshal(rw)
	if err != nil {
		return err
	}
	_, err = writer.Write(doc)
	return err
}

// AddParam ...
func (lr *LogoutRequest) AddParam(_, _ string) {

}
