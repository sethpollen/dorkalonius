// Basic game logic.

package dorkalonius

import (
	"sort"
)

type Game struct {
	TargetWord     string
	AvailableWords *WordList
}

const (
	numAvailableWords = 35

	// Manual tuning parameters. We use a high bias for the target word
	// in order to get something interesting. We use a much smaller bias
	// for the available words, since we want them to mostly reflect a
	// typical selection of words.
	targetWordBias    float64 = 3.04e-3
	availableWordBias float64 = 3.04e-6
)

func NewGame(wordSet *Index) (*Game, error) {
	// Always use the COCA set for generating the target word, as it's the only
	// set with accurate part-of-speech tagging.
	cocaWords, err := GetCocaIndex()
	if err != nil {
		return nil, err
	}

	target := cocaWords.SampleAdjective(
		1,
		SamplerConfig{int64(targetWordBias * float64(cocaWords.Leaves))})

	wordList := wordSet.Sample(
		numAvailableWords,
		SamplerConfig{int64(availableWordBias * float64(wordSet.Leaves))})

	sort.Sort(wordList)

	return &Game{target.Words[0].Word, wordList}, nil
}
