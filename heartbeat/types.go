package heartbeat

import "encoding/xml"

type (
	Response struct {
		XMLName xml.Name `xml:"res"`

		B body `xml:"body"`
	}

	body struct {
		Token string `xml:"token"`
	}
)
