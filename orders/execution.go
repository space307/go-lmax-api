package orders

import (
	"github.com/shopspring/decimal"
)

type (
	// Executions ...
	Executions struct {
		ExecutionId int64     `xml:"executionId"`
		Execution   Execution `xml:"execution"`
	}

	// Execution ...
	Execution struct {
		Price              decimal.Decimal `xml:"price"`
		Quantity           decimal.Decimal `xml:"quantity"`
		EncodedExecutionId string          `xml:"encodedExecutionId"`
	}
)
