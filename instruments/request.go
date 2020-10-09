package instruments

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/space307/go-lmax-api/model"

	"github.com/sirupsen/logrus"
)

const (
	requestURI = "/secure/instrument/searchCurrentInstruments?q=%s&offset=%d"
)

type (
	// InfoRequest ...
	InfoRequest struct {
		query  AssetClass
		offset int

		header map[string]string
	}
)

// NewInfoRequest ...
func NewInfoRequest(query AssetClass, offset int) *InfoRequest {
	return &InfoRequest{
		query:  query,
		offset: offset,
		header: make(map[string]string),
	}
}

// Header ...
func (r *InfoRequest) Header() map[string]string {
	return r.header
}

// RequestURI ...
func (r *InfoRequest) RequestURI() string {
	return fmt.Sprintf(requestURI, r.query, r.offset)
}

// Write ...
func (r *InfoRequest) Write(w io.Writer) error {
	return nil
}

// AddParam ...
func (r *InfoRequest) AddParam(key, value string) {
	r.header[key] = value
}

// GetInstrumentsInfo ...
func GetInstrumentsInfo(session model.Session, query AssetClass, onSuccess func(instruments []Info), onFailure func(code int, reader io.Reader)) error {
	go func() {
		var offset int
		var list []Info
		hasMore := true
		for hasMore {
			request := NewInfoRequest(query, offset)
			successCallback := func(reader io.Reader) {
				var response InfoList

				bytes, err := ioutil.ReadAll(reader)
				if err != nil {
					logrus.Error(err)
					return
				}
				if err := xml.Unmarshal(bytes, &response); err != nil {
					logrus.Error(err)
					return
				}

				list = append(list, response.Body.List.InfoList...)
				hasMore = response.Body.HasMoreResults
				offset = list[len(list)-1].ID
			}
			session.Get(request, successCallback, onFailure)
		}

		onSuccess(list)
	}()
	return nil
}
