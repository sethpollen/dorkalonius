// Defines an encoding/xml schema for storing inflection data as XML.

package wiktionary

import (
	"encoding/xml"
)

type Inflections struct {
  XMLName     xml.Name     `xml:"inflections"`
  Inflections []Inflection `xml:"inflection"`
}

type Inflection struct {
	XMLName        xml.Name `xml:"inflection"`
	BaseWord       string   `xml:"base"`
	Pos            string   `xml:"pos"`
	InflectedForms []string `xml:"inflected"`
}
