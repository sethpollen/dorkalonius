package dorkalonius_test

import (
	"strings"
	"testing"
)
import . "github.com/sethpollen/dorkalonius"

func TestBasic(t *testing.T) {
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
		"and again the quick brown fox jumps over the lazy dog", " ") {
		w.Add(WeightedWord{word, int64(i)})
		w.Check()
	}

	size = w.Size()
	if size != 10 {
		t.Errorf("Size: expected %d, got %d", 10, size)
	}
	weight = w.Weight()
	if weight != 55 {
		t.Errorf("Weight: expected %d, got %d", 55, weight)
	}

	words := w.GetWords()
	if len(words) != 10 {
		t.Errorf("len(words): expected %d, got %d", 10, len(words))
	}
	for i, expected := range []WeightedWord{
		WeightedWord{"dog", 10},
		WeightedWord{"the", 10},
		WeightedWord{"lazy", 9},
		WeightedWord{"over", 7},
		WeightedWord{"jumps", 6},
		WeightedWord{"fox", 5},
		WeightedWord{"brown", 4},
		WeightedWord{"quick", 3},
		WeightedWord{"again", 1},
		WeightedWord{"and", 0},
	} {
		if words[i] != expected {
			t.Errorf("words[%d]: expected %q, got %q", i, expected, words[i])
		}
	}
}

func TestAddAll(t *testing.T) {
	w1 := NewWordSet()
	w1.Add(WeightedWord{"a", 1})
	w1.Add(WeightedWord{"b", 3})
	w1.Add(WeightedWord{"c", 0})
	w1.Add(WeightedWord{"d", 15})

	w2 := NewWordSet()
	w2.Add(WeightedWord{"a", 15})
	w2.Add(WeightedWord{"b", 4})
	w2.Add(WeightedWord{"e", 1})
	w2.Add(WeightedWord{"f", 12})

	w1.AddAll(w2)
	w1.Check()

	expected := NewWordSet()
	expected.Add(WeightedWord{"a", 16})
	expected.Add(WeightedWord{"b", 7})
	expected.Add(WeightedWord{"c", 0})
	expected.Add(WeightedWord{"d", 15})
	expected.Add(WeightedWord{"e", 1})
	expected.Add(WeightedWord{"f", 12})

	if !wordSetsEqual(expected, w1) {
		t.Error()
	}
}

func wordSetsEqual(a, b WordSet) bool {
	aWords := a.GetWords()
	bWords := b.GetWords()

	if len(aWords) != len(bWords) {
		return false
	}

	for i, _ := range aWords {
		if aWords[i] != bWords[i] {
			return false
		}
	}

	return true
}
