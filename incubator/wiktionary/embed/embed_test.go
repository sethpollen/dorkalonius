package embed_test

import (
  "testing"
)
import . "github.com/sethpollen/dorkalonius/incubator/wiktionary/embed"

func TestSample(t *testing.T) {
  inflectionMap := GetInflectionMap()
  numBaseWords = inflectionMap.NumBaseWords()
  if numBaseWords != 1 {
    t.Errorf("Expected 1 base words; got %v", numBaseWords)
  }
}

// TODO: crew test