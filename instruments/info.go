package instruments

import (
	"encoding/xml"

	"github.com/shopspring/decimal"
)

const (
	// CurrencyAsset ...
	CurrencyClass = "CURRENCY"
	// IndexAsset ...
	IndexAsset = "INDEX"
)

type (
	// InfoList ...
	InfoList struct {
		XMLName xml.Name `xml:"res"`

		Body body `xml:"body"`
	}

	body struct {
		List           instruments `xml:"instruments"`
		HasMoreResults bool        `xml:"hasMoreResults"`
	}

	instruments struct {
		InfoList []Info `xml:"instrument"`
	}

	// Info ...
	Info struct {
		ID                             int             `xml:"id"`
		Name                           string          `xml:"name"`
		StartTime                      string          `xml:"startTime"`
		TradingHours                   TradingHours    `xml:"tradingHours"`
		Margin                         decimal.Decimal `xml:"margin"`
		Currency                       string          `xml:"currency"`
		UnitPrice                      decimal.Decimal `xml:"unitPrice"`
		MinimumOrderQuantity           decimal.Decimal `xml:"minimumOrderQuantity"`
		OrderQuantityIncrement         decimal.Decimal `xml:"orderQuantityIncrement"`
		MinimumPrice                   decimal.Decimal `xml:"minimumPrice"`
		MaximumPrice                   decimal.Decimal `xml:"maximumPrice"`
		TrustedSpread                  decimal.Decimal `xml:"trustedSpread"`
		PriceIncrement                 decimal.Decimal `xml:"priceIncrement"`
		StopBuffer                     decimal.Decimal `xml:"stopBuffer"`
		AssetClass                     AssetClass      `xml:"assetClass"`
		UnderlyingIsin                 string          `xml:"underlyingIsin"`
		Symbol                         string          `xml:"symbol"`
		MaximumPositionThreshold       decimal.Decimal `xml:"maximumPositionThreshold"`
		AggressiveCommissionRate       decimal.Decimal `xml:"aggressiveCommissionRate"`
		PassiveCommissionRate          decimal.Decimal `xml:"passiveCommissionRate"`
		MinimumCommission              decimal.Decimal `xml:"minimumCommission"`
		LongSwapPoints                 decimal.Decimal `xml:"longSwapPoints"`
		ShortSwapPoints                decimal.Decimal `xml:"shortSwapPoints"`
		DailyInterestRateBasis         decimal.Decimal `xml:"dailyInterestRateBasis"`
		ContractUnitOfMeasure          string          `xml:"contractUnitOfMeasure"`
		ContractSize                   decimal.Decimal `xml:"contractSize"`
		FundingBaseRate                string          `xml:"fundingBaseRate"`
		TradingDays                    TradingDays     `xml:"tradingDays"`
		RetailVolatilityBandPercentage decimal.Decimal `xml:"retailVolatilityBandPercentage"`
	}

	// AssetClass ...
	AssetClass string

	// TradingHours ...
	TradingHours struct {
		OpeningOffset int    `xml:"openingOffset"`
		ClosingOffset int    `xml:"closingOffset"`
		Timezone      string `xml:"timezone"`
	}

	// TradingDays ...
	TradingDays struct {
		TradingDays []string `xml:"tradingDay"`
	}
)
