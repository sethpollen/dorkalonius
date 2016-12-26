package wiktionary_test

import (
  "log"
	"testing"
)
import . "github.com/sethpollen/dorkalonius/wiktionary"

func loadMap() *InflectionMap {
  i, err := InflectionMapFromBzippedXml("./inflections.xml.bz2")
  if err != nil {
    log.Fatalln(err)
  }
  return i
}

func TestBasic(t *testing.T) {
  var inflectionMap = loadMap()
	if len(inflectionMap.BaseWords) != 222790 {
		t.Errorf("Got %d base words", len(inflectionMap.BaseWords))
	}
	if len(inflectionMap.InflectedToBase) != 242271 {
		t.Errorf("Got %d inflections", len(inflectionMap.InflectedToBase))
	}
  cases := [][]string{
    []string{"clothe", "clothe"},
    []string{"clothes", "clothe"},
    []string{"crew", "crew"},
    []string{"bear", "bear"},
    []string{"bears", "bear"},
    []string{"bearing", "bearing"},
    []string{"bearings", "bearing"},
  }
  for _, testCase := range cases {
    input := testCase[0]
    expected := testCase[1]
    actual := inflectionMap.GetBaseWord(input)
    if expected != actual {
      t.Errorf("Expected %q --> %q; got %q", input, expected, actual)
    }
  }
}