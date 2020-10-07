package account

import (
	"github.com/shopspring/decimal"
)

type (
	// Wallet ...
	Wallet struct {
		Currency string          `xml:"currency"`
		Balance  decimal.Decimal `xml:"balance"`
		Cash     decimal.Decimal `xml:"cash"`
	}

	// Wallets ...
	Wallets struct {
		Wallets []Wallet `xml:"wallet"`
	}
)
