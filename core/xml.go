package core

import "encoding/xml"

func unmarshalXML(d *xml.Decoder, start xml.StartElement, in interface{}) (int, int, error) {
	var offset, oldOffset int
	for {
		t, err := d.Token()
		if err != nil {
			return 0, 0, err
		}

		if se, ok := t.(xml.StartElement); ok && se.Name == start.Name {
			err = d.DecodeElement(&in, &se)
			if err != nil {
				return 0, 0, err
			}
			offset = int(d.InputOffset())
			if err := d.Skip(); err != nil {
				return oldOffset, offset, err
			}
			break
		}
		oldOffset = int(d.InputOffset())
	}
	return oldOffset, offset, nil
}
