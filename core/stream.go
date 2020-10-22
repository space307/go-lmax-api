package core

import (
	"errors"
	"io"
	"sync"

	"github.com/space307/go-lmax-api/model"
)

const (
	oneMB = 1024 * 1024

	eventsTag = "</events>"
)

type (
	// Decoder ...
	Decoder struct {
		r io.ReadCloser
	}
)

var EOFMsg = errors.New("eof")

// NewDecoder ...
func NewDecoder(r io.ReadCloser) *Decoder {
	return &Decoder{r: r}
}

// Close ...
func (d *Decoder) Close() error {
	return d.r.Close()
}

// Decode ...
func (d *Decoder) Decode() ([]byte, error) {
	const length = 256 * 1024
	chunk := make([]byte, length)
	n, err := d.r.Read(chunk)
	_ = n
	if err != nil && err != io.EOF {
		return nil, err
	}
	return chunk[:n], nil
}

// Stream ...
type Stream struct {
	model.Stream
	mutex   sync.Mutex
	source  *Decoder
	handler model.StreamHandler
	stopped bool
}

// NewStream ...
func NewStream(source *Decoder, handler model.StreamHandler) model.Stream {
	stream := &Stream{
		source:  source,
		handler: handler,
		mutex:   sync.Mutex{},
	}
	return stream
}

// Serve ...
func (s *Stream) Serve() error {
	return s.receive()
}

// Stop ...
func (s *Stream) Stop() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if !s.stopped {
		s.stopped = true
	}
}

func (s *Stream) isStopped() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.stopped
}

//receive reads result from io, decodes the data and sends it to the result channel
func (s *Stream) receive() error {

	buffer := make([]byte, 0, oneMB)
	lTag := len(eventsTag)

	defer s.source.Close()
	for {
		if s.isStopped() {
			break
		}

		var chunk []byte
		var err error
		for err == nil {
			chunk, err = s.source.Decode()
			buffer = append(buffer, chunk...)

			end := len(buffer)
			if end < lTag || string(buffer[end-lTag:end]) == eventsTag {
				err = io.EOF
			}
		}

		if err == io.EOF {
			s.handler.Handle(buffer)
		} else if err != nil {
			return err
		}

		buffer = buffer[:0]
	}
	return nil
}
