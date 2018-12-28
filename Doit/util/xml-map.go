package util

import (
	"encoding/xml"
	"io"
)

//XML转为字典
func XMLToMap(reader io.Reader, ignoreFirst bool) (map[string]string, error) {
	params := make(map[string]string)
	d := xml.NewDecoder(reader)
	value := ""
	for {
		token, err := d.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		switch t := token.(type) {
		case xml.StartElement:
			if ignoreFirst {
				ignoreFirst = false
				continue
			}
			value = t.Name.Local
		case xml.CharData:
			if value != "" {
				params[value] = string(t)
			}
		case xml.EndElement:
			value = ""
		}
	}
	return params, nil
}
