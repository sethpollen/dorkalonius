// Defines an encoding/xml schema for storing inflection data as XML.

package wiktionary

import (
	"encoding/xml"
)

type Inflection struct {
	XMLName        xml.Name `xml:"inflection"`
	BaseWord       string   `xml:"base"`
	Pos            string   `xml:"pos"`
	InflectedForms []string `xml:"inflected"`
}
