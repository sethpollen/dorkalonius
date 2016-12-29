// Basic game logic.

package dorkalonius

import (
  "sort"
)

type Game struct {
  TargetWord     string
  AvailableWords *WordList
}

const sampleSize = 35

func NewGame(sampler *Index) *Game {
  // Use a value of 1000000 here to get more interesting adjectives.
  adjective := sampler.SampleAdjective(1, SamplerConfig{1000000})

  // The least frequent words in our COCA corpus occur about 5000 times, so
  // this value of 1000 provides only a small boost to unlikely words.
  wordList := sampler.Sample(sampleSize, SamplerConfig{1000})
  sort.Sort(wordList)

  return &Game{adjective.Words[0].Word, wordList}
}
