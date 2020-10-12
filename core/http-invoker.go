package core

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/space307/go-lmax-api/model"

	"github.com/sirupsen/logrus"
)

const (
	defaultTimeout = time.Millisecond * 10000

	scheme = "http://"
)

type (
	// HttpCallback ...
	HttpCallback = func(code int, header model.Header, reader io.Reader) error

	// HttpInvoker ...
	HttpInvoker struct {
		addr string

		client       http.Client
		streamClient http.Client

		streamReader model.Stream
	}

	RequestParams map[string]string
)

// NewHttpInvoker ...
func NewHttpInvoker(addr string) *HttpInvoker {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}

	const size = 32768
	streamTr := &http.Transport{
		ReadBufferSize: size,
	}

	return &HttpInvoker{
		addr:         addr,
		client:       http.Client{Timeout: defaultTimeout, Transport: tr},
		streamClient: http.Client{Transport: streamTr},
	}
}

// Get ...
func (inv *HttpInvoker) Get(uri string, header map[string]string, args io.Reader, callback HttpCallback) error {
	req, err := http.NewRequest(http.MethodGet, scheme+inv.addr+uri, args)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "text/xml")
	req.Header.Set("Content-Type", "text/xml; UTF-8")

	for k, v := range header {
		req.Header.Add(k, v)
	}

	resp, err := inv.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if callback == nil {
		logrus.Error("httpInvoker: empty callback")
	}
	return callback(resp.StatusCode, resp.Header, resp.Body)
}

// Post ...
func (inv *HttpInvoker) Post(uri string, header map[string]string, args io.Reader, callback HttpCallback) error {
	req, err := http.NewRequest(http.MethodPost, scheme+inv.addr+uri, args)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "text/xml")
	req.Header.Set("Content-Type", "text/xml; UTF-8")

	for k, v := range header {
		req.Header.Add(k, v)
	}

	resp, err := inv.client.Do(req)
	if err != nil {
		return err
	}

	code := resp.StatusCode
	rHeader := resp.Header.Clone()

	data, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if err != nil {
		return err
	}

	if callback == nil {
		logrus.Error("httpInvoker: empty callback")
	}

	return callback(code, rHeader, bytes.NewBuffer(data))
}

// Post ...
func (inv *HttpInvoker) Stream(uri string, header map[string]string, args io.Reader, h model.StreamHandler) error {
	req, err := http.NewRequest(http.MethodPost, scheme+inv.addr+uri, args)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "text/xml")
	req.Header.Set("Content-Type", "text/xml; UTF-8")

	for k, v := range header {
		req.Header.Add(k, v)
	}

	resp, err := inv.streamClient.Do(req)
	if err != nil {
		return err
	}
	return inv.readStream(resp, h)
}

func (inv *HttpInvoker) readStream(resp *http.Response, h model.StreamHandler) error {
	dec := NewDecoder(resp.Body)
	inv.streamReader = NewStream(dec, h)
	return inv.streamReader.Serve()
}

// StopStreaming ...
func (inv *HttpInvoker) StopStreaming() {
	if inv.streamReader != nil {
		inv.streamReader.Stop()
	}
}
