package dorkalonius

import (
	"math/rand"
)

// Allows efficient weighted random sampling of a WordList. The index is
// stored as a binary tree, where the leaves correspond 1:1 woth Words in
// 'list'.
type Index struct {
	// The number of leaf nodes in the subtree rooted at this node.
	Leaves int
	// The total occurrences of all Words pointed to by nodes in the subtree
	// rooted at this node.
	Occurrences int64
	// Children of this node. If the node has only one child, 'right' will be
	// nil.
	Left  *Index
	Right *Index
	// Only set for leaf nodes, where 'left' and 'right' are both nil.
	Word *Word
}

// Parameters for word sampling.
type SamplerConfig struct {
  // Amount by which to increase the occurrences value for every word.
  // Higher values here make rare words more likely to be chosen.
  BaseOccurrences int64
}

func NewIndex(list *WordList) *Index {
	// Build the lowest level of the index tree.
	level := make([]*Index, list.Len())
	for i := range list.Words {
		word := &list.Words[i]
		level[i] = &Index{1, word.Occurrences, nil, nil, word}
	}

	// Build the internal levels of the tree.
	for len(level) > 1 {
		newLevel := make([]*Index, (len(level)+1)/2)
		for i := range newLevel {
			left := level[2*i]

			leaves := left.Leaves
			occurrences := left.Occurrences

			var right *Index = nil
			if 2*i+1 < len(level) {
				right = level[2*i+1]

				leaves += right.Leaves
				occurrences += right.Occurrences
			}

			newLevel[i] = &Index{leaves, occurrences, left, right, nil}
		}
		level = newLevel
	}

	return level[0]
}

func (self *Index) Weight(config SamplerConfig) int64 {
	return self.Occurrences + int64(self.Leaves) * config.BaseOccurrences
}

func (self *Index) LookUp(weight int64, config SamplerConfig) *Word {
	if self.Word != nil {
		return self.Word
	}
	leftWeight := self.Left.Weight(config)
	if weight < leftWeight {
		return self.Left.LookUp(weight, config)
	}
	return self.Right.LookUp(weight-leftWeight, config)
}

func (self *Index) pickWord(config SamplerConfig) *Word {
	return self.LookUp(rand.Int63n(self.Weight(config)), config)
}

func (self *Index) sampleMaybeAdjective(n int, config SamplerConfig,
                                        mustBeAdjective bool) *WordList {
  used := make(map[*Word]bool)
  result := NewWordList()
  for result.Len() < n {
    word := self.pickWord(config)
    if mustBeAdjective && !word.Adjective {
      continue
    }
    if used[word] {
      continue
    }
    used[word] = true
    result.AddWord(*word)
  }
  return result
}

// Randomly samples 'n' unique words.
func (self *Index) Sample(n int, config SamplerConfig) *WordList {
  return self.sampleMaybeAdjective(n, config, false)
}

// Randomly samples 'n' unique adjectives.
func (self *Index) SampleAdjective(n int, config SamplerConfig) *WordList {
  return self.sampleMaybeAdjective(n, config, true)
}
