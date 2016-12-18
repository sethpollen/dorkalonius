// Utilities for working with the Wiktionary XML dump.

package dump

import (
	"encoding/xml"
	"io"
)

// Result type of the XML parser; represents one MediaWiki page.
type Page struct {
	Title string `xml:"title"`
	Text  string `xml:"revision>text"`
}

type ProcessPage func(page *Page)

// Reads the full XML dump from 'file' and invokes 'process' on each page found
// in the dump.
func ReadDump(file io.Reader, process ProcessPage) {
	decoder := xml.NewDecoder(file)

	for {
		rawToken, _ := decoder.Token()
		if rawToken == nil {
			break
		}

		switch token := rawToken.(type) {
		case xml.StartElement:
			if token.Name.Local == "page" {
				var page Page
				decoder.DecodeElement(&page, &token)
				process(&page)
			}
		}
	}
}
