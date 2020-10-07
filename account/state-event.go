package account

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
		XMLName xml.Name `xml:"accountState"`

		AccountID               int64           `xml:"accountId"`
		Balance                 decimal.Decimal `xml:"balance"`
		Cash                    decimal.Decimal `xml:"cash"`
		AvailableFunds          decimal.Decimal `xml:"availableFunds"`
		AvailableToWithdraw     decimal.Decimal `xml:"availableToWithdraw"`
		UnrealisedProfitAndLoss decimal.Decimal `xml:"unrealisedProfitAndLoss"`
		Margin                  decimal.Decimal `xml:"margin"`

		Wallets Wallets `xml:"wallets"`

		Active bool `xml:"active"`
	}
)

// Type ...
func (_ *StateEvent) Type() events.Type {
	return events.AccountState
}

// SubscribeState ...
func SubscribeState(session model.Session, onSuccess func(reader io.Reader), onFailure func(code int, reader io.Reader)) error {
	subscription := model.NewSubscription()
	subscription.SetType(model.AccountStateSubscription)
	return session.Post(subscription, onSuccess, onFailure)
}
