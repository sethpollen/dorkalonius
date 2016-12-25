package wiktionary_test

import (
  "testing"
)
import . "github.com/sethpollen/dorkalonius/wiktionary"

func TestBasic(t *testing.T) {
  i, err := InflectionMapFromBzippedXml("./inflections.xml.bz2")
  if err != nil {
    t.Error(err)
    return
  }
  if len(i.BaseWords) != 0 {
    t.Errorf("Got %d base words", 0, len(i.BaseWords));
  }
  if len(i.InflectedToBase) != 0 {
    t.Errorf("Got %d inflections", 0, len(i.InflectedToBase));
  }
}