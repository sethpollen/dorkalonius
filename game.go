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
  targetWordBias float64 = 3.04e-3
  availableWordBias float64 = 3.04e-6
)

func NewGame(sampler *Index) *Game {
  target := sampler.SampleAdjective(
    1,
    SamplerConfig{int64(targetWordBias * float64(sampler.Leaves))})

  wordList := sampler.Sample(
    numAvailableWords,
    SamplerConfig{int64(availableWordBias * float64(sampler.Leaves))})

  sort.Sort(wordList)

  return &Game{target.Words[0].Word, wordList}
}
