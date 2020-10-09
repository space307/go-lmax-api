package main

import (
	"context"
	"fmt"

	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/space307/go-lmax-api"
	"github.com/space307/go-lmax-api/account"
	"github.com/space307/go-lmax-api/events"
	"github.com/space307/go-lmax-api/instruments"
	"github.com/space307/go-lmax-api/model"
	"github.com/space307/go-lmax-api/orders"
	"github.com/space307/go-lmax-api/positions"

	"github.com/go-xmlfmt/xmlfmt"
	"github.com/sirupsen/logrus"
)

type LoginClient struct {
}

func readToStr(reader io.Reader) string {
	buffer, err := ioutil.ReadAll(reader)
	if err != nil {
		logrus.Error(err)
	}
	return string(buffer)
}

func (lc *LoginClient) OnSuccess(session model.Session) {
	logrus.Infof("app: login success:\n%v", session)

	session.AddEventListener(events.AccountState, lc)
	session.AddEventListener(events.Positions, lc)
	session.AddEventListener(events.Orders, lc)
	session.AddEventListener(events.Position, lc)
	session.AddEventListener(events.Order, lc)

	if err := account.SubscribeState(session, func(reader io.Reader) {
		logrus.Infof("account subscription : \n%s", readToStr(reader))
	}, func(code int, reader io.Reader) {
		logrus.Errorf("account subscription : \n%s", readToStr(reader))
	}); err != nil {
		logrus.Error(err)
	}

	if err := positions.SubscribeState(session, func(_ io.Reader) {}, func(code int, reader io.Reader) {
		logrus.Errorf("position subscription : \n%s", readToStr(reader))
	}); err != nil {
		logrus.Error(err)
	}

	if err := orders.SubscribeState(session, func(_ io.Reader) {}, func(code int, reader io.Reader) {
		logrus.Errorf("orders subscription : \n%s", readToStr(reader))
	}); err != nil {
		logrus.Error(err)
	}

	stateRequest := func() error {
		if err := account.GetState(session, func(_ io.Reader) {}, func(code int, reader io.Reader) {
			logrus.Errorf("state : \n%s", xmlfmt.FormatXML(readToStr(reader), "\t", "    "))
		}); err != nil {
			return err
		}
		return nil
	}

	go func() {
		//for {
		time.Sleep(time.Second * 3)
		if err := stateRequest(); err != nil {
			logrus.Error(err)
		}
		//}
	}()
	go func() {
		if err := session.Serve(); err != nil {
			logrus.Error(err)
		}
	}()

	instruments.GetInstrumentsInfo(session, "", func(instInfo []instruments.Info) {
		logrus.Infof("instruments : \n%+v", instInfo)
	}, func(code int, reader io.Reader) {
		logrus.Errorf("failed to load instruments : code %d\n%s", code, xmlfmt.FormatXML(readToStr(reader), "", "\t"))
	})

	ctx := context.Background()
	<-ctx.Done()

	if err := session.Logout(account.NewLogoutRequest(), func(r io.Reader) {
		if err := session.Stop(); err != nil {
			logrus.Error(err)
		}
		logrus.Infof("logout : \n%s", readToStr(r))
	}, func(code int, r io.Reader) {
		if err := session.Stop(); err != nil {
			logrus.Error(err)
		}
		logrus.Errorf("logout : \n%s", readToStr(r))
	}); err != nil {
		logrus.Error(err)
	}
}

func (lc *LoginClient) OnFailure(code int, reader io.Reader) {
	logrus.Errorf("app: login failure <%d>:\n%s", code, readToStr(reader))
}

// OnEvent ...
func (lc *LoginClient) OnEvent(event events.Object) {
	switch event.Type() {
	case events.AccountState:
		logrus.Infof("event account: id - %d", event.(*account.StateEvent).AccountID)
	case events.Positions:
		logrus.Infof("event positions: len - %d", len(event.(*positions.StateEvent).Page.Positions))
	case events.Orders:
		logrus.Infof("event orders: len - %d", len(event.(*orders.StateEvent).Page.Orders))
	case events.Position:
		logrus.Infof("event position: id %d", event.(*positions.Event).Position.InstrumentID)
	case events.Order:
		logrus.Infof("event order: id %d", event.(*orders.Event).Order.InstrumentID)
	default:
		logrus.Warnf("app: unknown event %d", event.Type())
	}
}

func main() {
	if len(os.Args) != 5 {
		fmt.Println("Usage:", os.Args[0], "host", "username", "password", "scheme")
		return
	}

	host := os.Args[1]
	username := os.Args[2]
	password := os.Args[3]
	scheme := os.Args[4]

	if scheme != account.CfdDemo && scheme != account.CfdLive {
		fmt.Println("Invalid Scheme: must be ", account.CfdDemo, " | ", account.CfdLive)
		return
	}

	api := lmax.NewAPI(host)
	login := account.NewLoginRequest(username, password, account.CfdDemo)

	callback := &LoginClient{}
	if err := api.Login(login, callback); err != nil {
		logrus.Error(err)
	}
}
