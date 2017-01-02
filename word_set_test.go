package dorkalonius_test

import (
  "fmt"
  "strings"
  "testing"
)
import . "github.com/sethpollen/dorkalonius"

func TestBasicAdd(t *testing.T) {
  w := NewWordSet()
  w.Check()
  for i, word := range strings.Split(
      "the quick brown fox jumps over the lazy dog", " ") {
    w.Add(WeightedWord{word, int64(i)})
    w.Check()
  }
  fmt.Println(w.Size())
  fmt.Println(w.DebugString())
  t.Fail()
}