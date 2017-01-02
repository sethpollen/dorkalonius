package dorkalonius_test

import (
  "strings"
  "testing"
)
import . "github.com/sethpollen/dorkalonius"

func TestBasicAdd(t *testing.T) {
  w := NewWordSet()
  w.Check()
  
  size := w.Size()
  if size != 0 {
    t.Errorf("Size: expected %d, got %d", 0, size)
  }
  weight := w.Weight()
  if weight != 0 {
    t.Errorf("Weight: expected %d, got %d", 0, weight)
  }
  
  for i, word := range strings.Split(
      "the quick brown fox jumps over the lazy dog", " ") {
    w.Add(WeightedWord{word, int64(i)})
    w.Check()
  }
  
  size = w.Size()
  if size != 8 {
    t.Errorf("Size: expected %d, got %d", 8, size)
  }
  weight = w.Weight()
  if weight != 36 {
    t.Errorf("Weight: expected %d, got %d", 36, weight)
  }
}