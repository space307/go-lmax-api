package model

import "encoding/xml"

type (
	// Account ...
	Account struct {
		XMLName xml.Name `xml:"body"`

		Username                string `xml:"username"`
		Currency                string `xml:"currency"`
		AccountID               int64  `xml:"accountId"`
		RegistrationLegalEntity string `xml:"registrationLegalEntity"`
		DisplayLocale           string `xml:"displayLocale"`
		FundingDisallowed       bool   `xml:"fundingDisallowed"`
	}
)
