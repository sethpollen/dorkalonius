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
		w.Add(WeightedWord{word, int64(i + 1)})
		w.Check()
	}

	size = w.Size()
	if size != 10 {
		t.Errorf("Size: expected %d, got %d", 10, size)
	}
	weight = w.Weight()
	if weight != 66 {
		t.Errorf("Weight: expected %d, got %d", 66, weight)
	}

	words := w.GetWords()
	if len(words) != 10 {
		t.Errorf("len(words): expected %d, got %d", 10, len(words))
	}
	for i, expected := range []WeightedWord{
		WeightedWord{"the", 12},
		WeightedWord{"dog", 11},
		WeightedWord{"lazy", 10},
		WeightedWord{"over", 8},
		WeightedWord{"jumps", 7},
		WeightedWord{"fox", 6},
		WeightedWord{"brown", 5},
		WeightedWord{"quick", 4},
		WeightedWord{"again", 2},
		WeightedWord{"and", 1},
	} {
		if words[i] != expected {
			t.Errorf("words[%d]: expected %q, got %q", i, expected, words[i])
		}
	}
}

func TestInsert(t *testing.T) {
	w := NewWordSet()
	if !w.Insert(WeightedWord{"foo", 1}) {
		t.Error("Wrong Insert return value")
	}
	if !w.Insert(WeightedWord{"bar", 1}) {
		t.Error("Wrong Insert return value")
	}
	if w.Insert(WeightedWord{"foo", 1}) {
		t.Error("Wrong Insert return value")
	}
	if w.Size() != 2 {
		t.Error("Wrong size after insert")
	}
	if w.Weight() != 2 {
		t.Error("Wrong weight after insert")
	}
}

func TestAddAll(t *testing.T) {
	w1 := NewWordSet()
	w1.Add(WeightedWord{"a", 1})
	w1.Add(WeightedWord{"b", 3})
	w1.Add(WeightedWord{"c", 10})
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
	expected.Add(WeightedWord{"c", 10})
	expected.Add(WeightedWord{"d", 15})
	expected.Add(WeightedWord{"e", 1})
	expected.Add(WeightedWord{"f", 12})

	if !wordSetsEqual(expected, w1) {
		t.Error()
	}
}

func TestSample(t *testing.T) {
	w := NewWordSet()
	w.Add(WeightedWord{"a", 1})
	w.Add(WeightedWord{"d", 2})
	w.Add(WeightedWord{"f", 4})
	w.Add(WeightedWord{"b", 8})
	w.Add(WeightedWord{"e", 16})
	w.Add(WeightedWord{"c", 32})

	s := w.Sample(6, 0)
	if !wordSetsEqual(s, w) {
		t.Error()
	}

	for _, bias := range []int64{0, 10000} {
		counts := make(map[string]int)
		for i := 0; i < 10000; i++ {
			s = w.Sample(2, bias)
			if s.Size() != 2 {
				t.Error("Size: expected %d, got %d", 2, s.Size())
			}
			for _, word := range s.GetWords() {
				counts[word.Word]++
			}
		}

		// Do some rough probability checks.
		if counts["c"] < 5000 {
			t.Error("\"c\" should have been picked roughly 50% of the time")
		}
		if counts["e"] < 2500 {
			t.Error("\"e\" should have been picked roughly 25% of the time")
		}
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