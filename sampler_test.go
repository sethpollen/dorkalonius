package dorkalonius_test

import (
	"github.com/sethpollen/dorkalonius/coca"
	"testing"
)
import . "github.com/sethpollen/dorkalonius"

func TestSample(t *testing.T) {
	cocaWordList, err := coca.GetWordList()
	if err != nil {
		t.Error(err)
	}
	sampler := NewIndex(cocaWordList)
	var sampleSizes = []int{0, 1, 10, 100, 1000}
	for _, n := range sampleSizes {
		sample := sampler.Sample(n, SamplerConfig{10})
		if sample.Len() != n {
			t.Errorf("Expected sample of %v; got %v", n, sample.Len())
		}
	}
}
