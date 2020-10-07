package positions

import (
	"encoding/xml"
	"io"

	"github.com/space307/go-lmax-api/events"
	"github.com/space307/go-lmax-api/model"

	"github.com/shopspring/decimal"
)

type (
	// StateEvent ...
	StateEvent struct {
		XMLName xml.Name `xml:"positions"`

		Page Page `xml:"page"`
	}

	// Orders ...
	Page struct {
		Positions []Position `xml:"position"`
	}

	// Position ...
	Position struct {
		AccountID         int64           `xml:"accountId"`
		InstrumentID      int64           `xml:"instrumentId"`
		Valuation         decimal.Decimal `xml:"valuation"`
		ShortUnfilledCost decimal.Decimal `xml:"shortUnfilledCost"`
		LongUnfilledCost  decimal.Decimal `xml:"longUnfilledCost"`
		CumulativeCost    decimal.Decimal `xml:"cumulativeCost"`
		OpenQuantity      decimal.Decimal `xml:"openQuantity"`
		OpenCost          decimal.Decimal `xml:"openCost"`
	}
)

func (_ *StateEvent) Type() events.Type {
	return events.Positions
}

// SubscribeState ...
func SubscribeState(session model.Session, onSuccess func(reader io.Reader), onFailure func(code int, reader io.Reader)) error {
	subscription := model.NewSubscription()
	subscription.SetType(model.PositionsStateSubscription)
	return session.Post(subscription, onSuccess, onFailure)
}

func (s *StateEvent) Read(reader io.Reader) error {
	return nil
}
