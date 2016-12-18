package dorkalonius_test

import (
	"github.com/sethpollen/dorkalonius/coca"
	"testing"
)
import . "github.com/sethpollen/dorkalonius"

func TestSample(t *testing.T) {
	list := coca.GetWordList()
	sampler := NewIndex(list)
	var sampleSizes = []int{0, 1, 10, 100, 1000}
	for _, n := range sampleSizes {
		sample := sampler.Sample(n, 10)
		if sample.Len() != n {
			t.Errorf("Expected sample of %v; got %v", n, sample.Len())
		}
	}
}
