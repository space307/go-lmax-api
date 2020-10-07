package events

import (
	"encoding/xml"
)

const (
	// Unknown ...
	Unknown Type = -1

	// AccountState ...
	AccountState Type = 0

	// Positions ...
	Positions Type = 1

	// Orders ...
	Orders Type = 2

	// Positions ...
	Position Type = 3

	// Orders ...
	Order Type = 4
)

var (
	types = []string{
		"accountState",
		"positions",
		"orders",
		"position",
		"order",
	}
)

type (
	// Type ...
	Type = int
	// Object ...
	Object interface {
		Type() Type
	}
	// Observer ...
	Observer interface {
		OnEvent(event Object)
	}
	TypeExtractor struct {
		keys []string
	}
)

// GetType ...
func GetType(s string) Type {
	for i, str := range types {
		if s == str {
			return i
		}
	}
	return Unknown
}

// UnmarshalXML ...
func (te *TypeExtractor) UnmarshalXML(d *xml.Decoder) error {
	te.keys = make([]string, 0, 16)

	var readKey bool
	var eventKey string
	t, _ := d.Token()
	for t != nil {
		switch tt := t.(type) {
		case xml.StartElement:
			if tt.Name.Local == "body" {
				readKey = true
				break
			}
			if readKey {
				te.keys = append(te.keys, tt.Name.Local)
				eventKey = tt.Name.Local
				readKey = false
				break
			}
		case xml.EndElement:
			if eventKey == tt.Name.Local {
				eventKey = ""
				readKey = true
				break
			}
			if tt.Name.Local == "body" {
				readKey = false
			}
		}
		t, _ = d.Token()
	}
	return nil
}

// RawTypes ...
func (te *TypeExtractor) RawTypes() []string {
	return te.keys
}
