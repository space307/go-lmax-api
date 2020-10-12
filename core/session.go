package core

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/space307/go-lmax-api/account"
	"github.com/space307/go-lmax-api/events"
	"github.com/space307/go-lmax-api/heartbeat"
	"github.com/space307/go-lmax-api/model"
	"github.com/space307/go-lmax-api/orders"
	"github.com/space307/go-lmax-api/positions"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

const (
	userAgentParam = "User-Agent"
	cookieParam    = "cookie"

	inessensialErr = "unexpected end element"
)

type (
	// Session ...
	Session struct {
		XMLName xml.Name `xml:"res"`

		AccountDetails model.Account

		httpInvoker *HttpInvoker

		id        string
		userAgent string

		stream *stream
		en     *eventsNotifier
	}

	eventsNotifier struct {
		sync.RWMutex

		observers map[events.Type][]events.Observer
	}

	stream struct {
		sync.WaitGroup
		stopper chan struct{}
	}
)

func newEventsNotifier() *eventsNotifier {
	return &eventsNotifier{observers: make(map[events.Type][]events.Observer)}
}

func (en *eventsNotifier) Handle(raw []byte) {
	var te events.TypeExtractor
	if err := te.UnmarshalXML(xml.NewDecoder(bytes.NewBuffer(raw))); err != nil {
		logrus.Error(err)
		return
	}

	if len(te.RawTypes()) > 0 {
		eventsArr := make([]events.Object, 0, len(te.RawTypes()))
		buffer := make([]byte, 0, len(raw))
		buffer = append(buffer, raw...)

		var start xml.StartElement
		var event events.Object
		for _, key := range te.RawTypes() {
			switch events.GetType(key) {
			case events.AccountState:
				event = &account.StateEvent{}
			case events.Positions:
				event = &positions.StateEvent{}
			case events.Orders:
				event = &orders.StateEvent{}
			case events.Position:
				event = &positions.Event{}
			case events.Order:
				event = &orders.Event{}
			case events.Heartbeat:
				event = &heartbeat.Event{}
			default:
				logrus.Warnf("unknown event type <%s>", key)
				continue
			}
			start = xml.StartElement{Name: xml.Name{Local: key}}
			decoder := xml.NewDecoder(bytes.NewBuffer(buffer))
			offset, end, err := unmarshalXML(decoder, start, event)
			if err != nil && !strings.Contains(err.Error(), inessensialErr) {
				logrus.Error(err)
				break
			}
			length := len(buffer)
			copy(buffer[offset:], buffer[end:])
			buffer = buffer[:length-end+offset]

			eventsArr = append(eventsArr, event)
		}

		for _, e := range eventsArr {
			if observers, found := en.observers[e.Type()]; found {
				for _, o := range observers {
					o.OnEvent(e)
				}
			}
		}
	}
}

// Session ...
func (s *Session) String() string {
	return fmt.Sprintf("id: %s\nuserAgent: %s\naccount: %v\n", s.id, s.userAgent, s.AccountDetails)
}

// NewSession ...
func NewSession(header model.Header, reader io.Reader, invoker *HttpInvoker) (model.Session, error) {
	s := &Session{
		httpInvoker: invoker,
		en:          newEventsNotifier(),
		stream:      &stream{stopper: make(chan struct{})},
	}

	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	if err := xml.Unmarshal(bytes, s); err != nil {
		return nil, err
	}
	s.userAgent = createUserAgent(s.AccountDetails.Username)

	cookie := extractCookie(header)
	s.id = cookie.sessionID
	return s, nil
}

// Logout ...
func (s *Session) Logout(request model.Request, success func(r io.Reader), failure func(code int, r io.Reader)) error {
	return s.Post(request, success, failure)
}

// Get ...
func (s *Session) Get(request model.Request, success func(r io.Reader), failure func(code int, r io.Reader)) error {
	args := bytes.NewBuffer(nil)
	err := request.Write(args)
	if err != nil {
		return err
	}

	request.AddParam(cookieParam, s.ID())
	request.AddParam(userAgentParam, s.UserAgent())

	return s.httpInvoker.Get(request.RequestURI(), request.Header(), args, func(code int, header model.Header, reader io.Reader) error {
		if code != http.StatusOK {
			failure(code, reader)
			return nil
		}
		success(reader)
		return nil
	})
}

// Post ...
func (s *Session) Post(request model.Request, success func(r io.Reader), failure func(code int, r io.Reader)) error {
	args := bytes.NewBuffer(nil)
	err := request.Write(args)
	if err != nil {
		return err
	}

	request.AddParam(cookieParam, s.ID())
	request.AddParam(userAgentParam, s.UserAgent())

	return s.httpInvoker.Post(request.RequestURI(), request.Header(), args, func(code int, header model.Header, reader io.Reader) error {
		if code != http.StatusOK {
			failure(code, reader)
			return nil
		}
		success(reader)
		return nil
	})
}

// Stream ...
func (s *Session) Stream(request model.Request, h model.StreamHandler) error {
	args := bytes.NewBuffer(nil)
	err := request.Write(args)
	if err != nil {
		return err
	}

	request.AddParam(cookieParam, s.ID())
	request.AddParam(userAgentParam, s.UserAgent())

	return s.httpInvoker.Stream(request.RequestURI(), request.Header(), args, h)
}

func (s *Session) HeartbeatRequest(request model.Request, success func(), failure func(code int, r io.Reader), disconnected func()) error {
	err := s.Post(request, func(reader io.Reader) {
		bytes, err := ioutil.ReadAll(reader)
		if err != nil {
			logrus.Error(err)
			return
		}
		var response heartbeat.Response
		if err := xml.Unmarshal(bytes, &response); err != nil {
			logrus.Error(err)
			return
		}
		success()
	}, func(code int, r io.Reader) {
		if code == http.StatusForbidden {
			disconnected()
		} else {
			failure(code, r)
		}
	})
	return err
}

// Serve ...
func (s *Session) Serve() error {
	return s.eventLoop(NewStreamRequest())
}

// Stop ...
func (s *Session) Stop() error {
	s.httpInvoker.StopStreaming()
	close(s.stream.stopper)
	return nil
}

// Wait ...
func (s *Session) Wait() {
	s.stream.Wait()
}

func (s *Session) eventLoop(request model.Request) error {
	s.stream.Add(1)
	defer s.stream.Done()

	for {
		if err := s.Stream(request, s.en); err != nil && err != io.EOF {
			return err
		}
		select {
		case <-s.stream.stopper:
			return nil
		default:
		}
	}
}

// AddEventListener ...
func (s *Session) AddEventListener(t events.Type, o events.Observer) {
	s.en.Lock()
	s.en.observers[t] = append(s.en.observers[t], o)
	s.en.Unlock()
}

// RemoveEventListener ...
func (s *Session) RemoveEventListener(t events.Type, o events.Observer) {
	s.en.Lock()
	if observers, found := s.en.observers[t]; found {
		for i, observer := range observers {
			if observer == o {
				observers[i] = observers[len(observers)-1]
				observers = observers[:len(observers)-1]
				break
			}
		}
		s.en.observers[t] = observers
	}
	s.en.Unlock()
}

// ID ...
func (s *Session) ID() string {
	return s.id
}

// UserAgent ...
func (s *Session) UserAgent() string {
	return s.userAgent
}
