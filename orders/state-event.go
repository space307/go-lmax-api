package orders

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
		XMLName xml.Name `xml:"orders"`

		Page Page `xml:"page"`
	}

	// Order ...
	Page struct {
		Orders []Order `xml:"order"`
	}

	// Order ...
	Order struct {
		TimeInForce           string          `xml:"timeInForce"`
		InstructionId         int64           `xml:"instructionId"`
		OriginalInstructionId int64           `xml:"originalInstructionId"`
		OrderId               string          `xml:"orderId"`
		AccountId             int64           `xml:"accountId"`
		InstrumentID          int64           `xml:"instrumentId"`
		Quantity              decimal.Decimal `xml:"quantity"`
		MatchedQuantity       decimal.Decimal `xml:"matchedQuantity"`
		MatchedCost           decimal.Decimal `xml:"matchedCost"`
		CancelledQuantity     decimal.Decimal `xml:"cancelledQuantity"`
		Timestamp             string          `xml:"timestamp"`
		OrderType             string          `xml:"orderType"`
		OpenQuantity          decimal.Decimal `xml:"openQuantity"`
		OpenCost              decimal.Decimal `xml:"openCost"`
		CumulativeCost        decimal.Decimal `xml:"cumulativeCost"`
		Commission            decimal.Decimal `xml:"commission"`
		StopReferencePrice    decimal.Decimal `xml:"stopReferencePrice"`
		StopLossOffset        string          `xml:"stopLossOffset"`
		StopProfitOffset      string          `xml:"stopProfitOffset"`
		WorkingState          string          `xml:"workingState"`

		Executions Executions `xml:"executions"`
	}
)

func (_ *StateEvent) Type() events.Type {
	return events.Orders
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
