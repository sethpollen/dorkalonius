// Defines a container for words which tracks word weights and allows random
// sampling. Words are stored in a sorted AVL tree, with each internal
// node storing the number and total weight of its descendant leaves.

package dorkalonius

import (
	"bytes"
	"fmt"
	"log"
	"sort"
	"strings"
)

type WeightedWord struct {
	Word   string
	Weight int64
}

type node struct {
	// Child pointers.
	Left  *node
	Right *node

	// The Weighted represented by this node.
	Word WeightedWord

	// Information about the subtree rooted at this node, which includes
	// this node.
	SubtreeHeight int
	SubtreeSize   int64
	SubtreeWeight int64
}

func newLeafNode(word WeightedWord) *node {
	// Newly inserted nodes start out as red.
	return &node{nil, nil, word, 1, 1, word.Weight}
}

type WordSet struct {
	// nil for an empty WordSet.
	root *node
}

func NewWordSet() WordSet {
	return WordSet{nil}
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC API

// Checks invariants and crashes if any are found.
func (self WordSet) Check() {
	if self.root == nil {
		return
	}
	self.check(self.root)

	// Check overall sortedness.
	var last string = ""
	visit(self.root, 0, func(n *node, depth int) {
		if last != "" && strings.Compare(last, n.Word.Word) >= 0 {
			log.Fatal("Not sorted")
		}
		last = n.Word.Word
	})
}

func (self WordSet) DebugString() string {
	// Add 1 to the height so we can print unbalanced trees.
	rows := make([]bytes.Buffer, subtreeHeight(self.root)+1)
	visit(self.root, 0, func(n *node, depth int) {
		if _, err := rows[depth].WriteString(
			fmt.Sprintf("%s:%d ", n.Word.Word, children(n))); err != nil {
			log.Fatal(err)
		}
	})

	var all bytes.Buffer
	for _, row := range rows {
		if row.Len() == 0 {
			break
		}
		if _, err := row.WriteTo(&all); err != nil {
			log.Fatal(err)
		}
		if _, err := all.WriteString("\n"); err != nil {
			log.Fatal(err)
		}
	}
	return all.String()
}

func (self WordSet) Size() int64 {
	return subtreeSize(self.root)
}

func (self WordSet) Weight() int64 {
	return subtreeWeight(self.root)
}

func (self *WordSet) Add(word WeightedWord) {
	if self.root == nil {
		self.root = newLeafNode(word)
		return
	}

	path := []*node{self.root}
	for {
		cur := path[len(path)-1]

		c := strings.Compare(word.Word, cur.Word.Word)
		if c == 0 {
			cur.Word.Weight += word.Weight
			// We didn't actually insert any nodes, but we break to the rebalance
			// call anyway in order to update subtree counts.
			break
		}

		if c < 0 {
			if cur.Left == nil {
				cur.Left = newLeafNode(word)
				path = append(path, cur.Left)
				break
			}
			path = append(path, cur.Left)
			continue
		}

		// c > 0
		if cur.Right == nil {
			cur.Right = newLeafNode(word)
			path = append(path, cur.Right)
			break
		}
		path = append(path, cur.Right)
	}

	self.rebalance(path)
}

func (self *WordSet) AddAll(other WordSet) {
	visit(other.root, 0, func(n *node, depth int) {
		self.Add(n.Word)
	})
}

// Gets the contents of this WordSet, sorted by descending weight.
func (self WordSet) GetWords() []WeightedWord {
	words := make([]WeightedWord, self.Size())
	i := 0
	visit(self.root, 0, func(n *node, depth int) {
		words[i] = n.Word
		i++
	})
	sort.Sort(SortWeightedWords(words))
	return words
}

///////////////////////////////////////////////////////////////////////////////
// HELPERS

// Returns the height of the subtree rooted at 'n'.
func (self *WordSet) check(n *node) {
	if n == nil {
		return
	}

	if n.Left != nil {
		if strings.Compare(n.Left.Word.Word, n.Word.Word) >= 0 {
			log.Fatal("Not ordered")
		}
	}
	if n.Right != nil {
		if strings.Compare(n.Word.Word, n.Right.Word.Word) >= 0 {
			log.Fatal("Not ordered")
		}
	}

	expectedSubtreeHeight :=
		max(subtreeHeight(n.Left), subtreeHeight(n.Right)) + 1
	if n.SubtreeHeight != expectedSubtreeHeight {
		log.Fatalf("Bad SubtreeHeight for %q: Expected %d, got %d",
			n.Word.Word, expectedSubtreeHeight, n.SubtreeHeight)
	}

	expectedSubtreeSize :=
		subtreeSize(n.Left) + subtreeSize(n.Right) + 1
	if n.SubtreeSize != expectedSubtreeSize {
		log.Fatalf("Bad SubtreeSize for %q: Expected %d, got %d",
			n.Word.Word, expectedSubtreeSize, n.SubtreeSize)
	}

	expectedSubtreeWeight :=
		subtreeWeight(n.Left) + subtreeWeight(n.Right) + n.Word.Weight
	if n.SubtreeWeight != expectedSubtreeWeight {
		log.Fatalf("Bad SubtreeWeight for %q: Expected %d, got %d",
			n.Word.Word, expectedSubtreeWeight, n.SubtreeWeight)
	}

	imb := imbalance(n)
	if abs(imb) > 1 {
		log.Fatalf("Too much imbalance (%d):\n%s", imb, self.DebugString())
	}

	self.check(n.Left)
	self.check(n.Right)
}

func updateSubtreeInfo(n *node) {
	n.SubtreeHeight = max(subtreeHeight(n.Left), subtreeHeight(n.Right)) + 1
	n.SubtreeSize = subtreeSize(n.Left) + subtreeSize(n.Right) + 1
	n.SubtreeWeight = subtreeWeight(n.Left) + subtreeWeight(n.Right) +
		n.Word.Weight
}

func (self *WordSet) rebalance(path []*node) {
	for i := len(path) - 1; i >= 0; i-- {
		n := path[i]

		// Propagate any changes from the last iteration.
		updateSubtreeInfo(n)

		imb := imbalance(n)
		if abs(imb) > 2 {
			log.Fatalf("Imbalance too large: %d", abs(i))
		}

		if imb == 0 {
			// No imbalance.
			continue
		}
		if i == 0 {
			if imb == 1 || imb == -1 {
				// We are OK with +-1 imbalance at the root.
				continue
			}
		} else {
			parent := path[i-1]
			if n == parent.Left {
				if imb == -1 {
					// We are OK with a heaver left subtree for a node which is on the
					// left side of its parent.
					continue
				}
			} else {
				if imb == 1 {
					// We are OK with a heaver right subtree for a node which is on
					// the right side of its parent.
					continue
				}
			}
		}

		// Find the pointer used by n's parent to refer to n.
		var parentPtr **node
		if i == 0 {
			parentPtr = &self.root
		} else {
			parent := path[i-1]
			if parent.Left == n {
				parentPtr = &parent.Left
			} else {
				parentPtr = &parent.Right
			}
		}

		// Perform the rotation.
		var child *node
		if imb < 0 {
			// The left subtree is too tall.
			child = n.Left
			n.Left = child.Right
			child.Right = n
		} else {
			// The right subtree is too tall.
			child = n.Right
			n.Right = child.Left
			child.Left = n
		}
		*parentPtr = child

		updateSubtreeInfo(n)
		updateSubtreeInfo(child)
	}
}

func visit(n *node, depth int, visitor func(n *node, depth int)) {
	if n == nil {
		return
	}
	visit(n.Left, depth+1, visitor)
	visitor(n, depth)
	visit(n.Right, depth+1, visitor)
}

// Returns a negative number if the left subtree is taller and a positive
// number if the right subtree is taller.
func imbalance(n *node) int {
	return -subtreeHeight(n.Left) + subtreeHeight(n.Right)
}

func subtreeHeight(n *node) int {
	if n == nil {
		return 0
	}
	return n.SubtreeHeight
}

func subtreeSize(n *node) int64 {
	if n == nil {
		return 0
	}
	return n.SubtreeSize
}

func subtreeWeight(n *node) int64 {
	if n == nil {
		return 0
	}
	return n.SubtreeWeight
}

func children(n *node) int {
	c := 0
	if n.Left != nil {
		c++
	}
	if n.Right != nil {
		c++
	}
	return c
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Support sorting of WeightedWords.

type SortWeightedWords []WeightedWord

func (self SortWeightedWords) Len() int {
	return len(self)
}

func (self SortWeightedWords) Less(i, j int) bool {
	if self[i].Weight != self[j].Weight {
		return self[i].Weight > self[j].Weight
	}
	return strings.Compare(self[i].Word, self[j].Word) < 0
}

func (self SortWeightedWords) Swap(i, j int) {
	temp := self[i]
	self[i] = self[j]
	self[j] = temp
}
