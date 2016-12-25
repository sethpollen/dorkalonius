// TODO: remove me once we fully build out choices.txt

package main

import (
  "github.com/sethpollen/dorkalonius/wiktionary"
  "log"
)

func main() {
  _, err := wiktionary.InflectionMapFromBzippedXml("./inflections.xml.bz2")
  if err != nil {
    log.Fatalln(err)
  }
}