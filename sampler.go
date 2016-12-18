package words

import (
	"math/rand"
	"strings"
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

func (self *Index) Weight(baseOccurrences int64) int64 {
	return self.Occurrences + int64(self.Leaves)*baseOccurrences
}

func (self *Index) LookUp(weight int64, baseOccurrences int64) *Word {
	if self.Word != nil {
		return self.Word
	}
	leftWeight := self.Left.Weight(baseOccurrences)
	if weight < leftWeight {
		return self.Left.LookUp(weight, baseOccurrences)
	}
	return self.Right.LookUp(weight-leftWeight, baseOccurrences)
}

func (self *Index) pickWord(baseOccurrences int64) *Word {
	return self.LookUp(rand.Int63n(self.Weight(baseOccurrences)),
		baseOccurrences)
}

// Randomly samples 'n' unique words the given 'partOfSpeech'. If
// 'partOfSpeech' is '*', then any part of speech may be returned.
func (self *Index) SamplePartOfSpeech(n int, partOfSpeech byte,
	baseOccurrences int64) *WordList {
	used := make(map[*Word]bool)
	result := NewWordList()
	for result.Len() < n {
		word := self.pickWord(baseOccurrences)
		if partOfSpeech != '*' &&
			strings.IndexByte(word.PartsOfSpeech, partOfSpeech) < 0 {
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
func (self *Index) Sample(n int, baseOccurrences int64) *WordList {
	return self.SamplePartOfSpeech(n, '*', baseOccurrences)
}
