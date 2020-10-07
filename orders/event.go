package orders

import (
	"encoding/xml"

	"github.com/space307/go-lmax-api/events"
)

type (
	// Event ...
	Event struct {
		XMLName xml.Name `xml:"order"`
		Order
	}
)

// Type ...
func (_ *Event) Type() events.Type {
	return events.Order
}
