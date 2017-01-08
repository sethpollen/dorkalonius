package dorkalonius_test

import (
	"bytes"
	"strings"
	"testing"
)
import . "github.com/sethpollen/dorkalonius"

func TestBasic(t *testing.T) {
	w := NewWordSet()
	if err := w.Check(); err != nil {
		t.Error(err)
		return
	}

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

		if err := w.Check(); err != nil {
			t.Error(err)
			return
		}

		var buf bytes.Buffer
		if err := w.Serialize(&buf); err != nil {
			t.Error(err)
			return
		}
		if buf.Len() == 0 {
			t.Error("Serialize produced empty buffer")
			return
		}
		deserialized, err := DeserializeWordSet(&buf)
		if err != nil {
			t.Error(err)
		}
		if !wordSetsEqual(w, *deserialized) {
			t.Error("Serialization/deserialization not faithful")
		}
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
	if err := w1.Check(); err != nil {
		t.Error(err)
		return
	}

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
				return
			}
			for _, word := range s.GetWords() {
				counts[word.Word]++
			}
		}

		if len(counts) < 6 {
			t.Error("Not all words were sampled")
		}

		// Do some rough probability checks.
		expectedCountC := 10000
		expectedCountE := 5000
		if bias > 0 {
			expectedCountC = 20000 / 6
			expectedCountE = 20000 / 6
		}

		if counts["c"] < expectedCountC/2 {
			t.Error("Count for \"c\" too small")
		}
		if counts["c"] > expectedCountC*2 {
			t.Error("Count for \"c\" too large")
		}
		if counts["e"] < expectedCountE/2 {
			t.Error("Count for \"e\" too small")
		}
		if counts["e"] > expectedCountE*2 {
			t.Error("Count for \"e\" too large")
		}
	}
}

func TestPrettyPrint(t *testing.T) {
	w := NewWordSet()
	for _, word := range strings.Split(
		"the quick brown fox", " ") {
		w.Add(WeightedWord{word, 1})
	}

	actual := strings.TrimSpace(w.PrettyPrint())
	expected := strings.TrimSpace(`
+-> quick
  +-> fox
  | +-> brown
  +-> the`)

	if actual != expected {
		t.Errorf("Actual:\n%s\nExpected:\n%s", actual, expected)
	}
}

// Helpers.

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
