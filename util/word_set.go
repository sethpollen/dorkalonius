// Defines a container for words which tracks word weights and allows random
// sampling. Words are stored in a sorted AVL tree, with each internal
// node storing the number and total weight of its descendant leaves.

package util

import (
  "bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"sort"
	"strings"
)

type WeightedWord struct {
	Word      string
	Weight    int64
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

// Checks invariants and returns an error if any are found.
func (self WordSet) Check() error {
	if err := self.check(self.root); err != nil {
		return err
	}

	// Check overall sortedness.
	var last string = ""
	var err error = nil
	visit(self.root, 0, func(n *node, depth int) {
		if err != nil {
			return
		}
		if last != "" && strings.Compare(last, n.Word.Word) >= 0 {
			err = errors.New("Not ordered")
		}
		last = n.Word.Word
	})
	return err
}

func (self WordSet) Size() int64 {
	return subtreeSize(self.root)
}

func (self WordSet) Weight() int64 {
	return subtreeWeight(self.root)
}

func (self *WordSet) Add(word WeightedWord) {
	self.add(word, false)
}

// Like Add, but requires that the word not already be present in the set.
// Returns true iff the insertion happened.
func (self *WordSet) Insert(word WeightedWord) bool {
	return self.add(word, true)
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

// Randomly samples 'n' words from this WordSet (using their Weights) and returns
// those words. 'nodeBias' will be added to every node's weight.
func (self WordSet) Sample(n int64, nodeBias int64) WordSet {
	if n > self.Size() {
		log.Fatalf("Cannot sample %d words from a WordSet of size %d",
			n, self.Size())
	}
	totalWeight := self.Weight() + nodeBias*self.Size()

	sample := NewWordSet()
	for sample.Size() < n {
		point := rand.Int63n(totalWeight)
		cur := self.root
		for {
			leftWeight := subtreeWeight(cur.Left) + nodeBias*subtreeSize(cur.Left)
			if point < leftWeight {
				cur = cur.Left
				continue
			}
			point -= leftWeight

			curWeight := cur.Word.Weight + nodeBias
			if point < curWeight {
				sample.Insert(cur.Word)
				break
			}
			point -= curWeight

			cur = cur.Right
		}
	}
	return sample
}

func (self WordSet) PrettyPrint() string {
  return prettyPrint(self.root)
}

///////////////////////////////////////////////////////////////////////////////
// SERIALIZATION

var byteOrder = binary.LittleEndian

func (self WordSet) Serialize(out io.Writer) error {
	return serialize(out, self.root)
}

func DeserializeWordSet(in io.Reader) (*WordSet, error) {
	root, err := deserialize(in)
	if err != nil {
		return nil, err
	}
	words := &WordSet{root}
	if err = words.Check(); err != nil {
		return nil, err
	}
	return words, nil
}

///////////////////////////////////////////////////////////////////////////////
// HELPERS

// If 'requireInsert' is true and the given word is already in this set, this
// method does nothing and returns false. Otherwise returns true.
func (self *WordSet) add(word WeightedWord, requireInsert bool) bool {
	if word.Weight <= 0 {
		log.Fatal("Weights must be positive")
	}

	if self.root == nil {
		self.root = newLeafNode(word)
		return true
	}

	path := []*node{self.root}
	for {
		cur := path[len(path)-1]

		c := strings.Compare(word.Word, cur.Word.Word)
		if c == 0 {
			if requireInsert {
				return false
			}
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
	return true
}

// Checks tree invariants. Returns an error on failure.
func (self *WordSet) check(n *node) error {
	if n == nil {
		return nil
	}

	if n.Word.Weight <= 0 {
		return fmt.Errorf("Nonpositive weight: %d", n.Word.Weight)
	}

	if n.Left != nil {
		if strings.Compare(n.Left.Word.Word, n.Word.Word) >= 0 {
			return errors.New("Not ordered")
		}
	}
	if n.Right != nil {
		if strings.Compare(n.Word.Word, n.Right.Word.Word) >= 0 {
			return errors.New("Not ordered")
		}
	}

	expectedSubtreeHeight :=
		max(subtreeHeight(n.Left), subtreeHeight(n.Right)) + 1
	if n.SubtreeHeight != expectedSubtreeHeight {
		return fmt.Errorf("Bad SubtreeHeight for %q: Expected %d, got %d",
			n.Word.Word, expectedSubtreeHeight, n.SubtreeHeight)
	}

	expectedSubtreeSize :=
		subtreeSize(n.Left) + subtreeSize(n.Right) + 1
	if n.SubtreeSize != expectedSubtreeSize {
		return fmt.Errorf("Bad SubtreeSize for %q: Expected %d, got %d",
			n.Word.Word, expectedSubtreeSize, n.SubtreeSize)
	}

	expectedSubtreeWeight :=
		subtreeWeight(n.Left) + subtreeWeight(n.Right) + n.Word.Weight
	if n.SubtreeWeight != expectedSubtreeWeight {
		return fmt.Errorf("Bad SubtreeWeight for %q: Expected %d, got %d",
			n.Word.Word, expectedSubtreeWeight, n.SubtreeWeight)
	}

	imb := imbalance(n)
	if abs(imb) > 1 {
		return fmt.Errorf("Too much imbalance (%d) for %q\n%s",
                      imb, n.Word.Word, prettyPrint(n))
	}

	if err := self.check(n.Left); err != nil {
		return err
	}
	if err := self.check(n.Right); err != nil {
		return err
	}
	return nil
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
			log.Fatalf("Imbalance (%d) too large for %q\n%s",
                 abs(i), n.Word.Word, self.PrettyPrint())
		}
		if abs(imb) < 2 {
			// Acceptable imbalance.
			continue
		}

		// Find the pointer used by n's parent to refer to n.
		var nPtr **node
		if i == 0 {
			nPtr = &self.root
		} else {
			parent := path[i-1]
			if parent.Left == n {
				nPtr = &parent.Left
			} else {
				nPtr = &parent.Right
			}
		}
    
    // Track the set of nodes which need their subtree info recalculated.
    recalculate := make([]*node, 0)
    
    if imb < 0 {
      // n.Left is too tall.
      if imbalance(n.Left) < 0 {
        // We must do an single right rotation.
        //     n             p
        //    / \           / \
        //   p   C   ==>   A   n
        //  / \               / \
        // A   B             B   C
        // (Where A is taller than B).
        p := n.Left
        
        *nPtr = p
        n.Left = p.Right
        p.Right = n
        
        recalculate = append(recalculate, p)
      } else {
        // We must do a left-right rotation.
        //     n                n              t
        //    / \              / \            / \
        //   p   D   ==>      t   D   ==>   p     n
        //  / \              / \           / \   / \
        // A   t            p   C         A   B C   D
        //    / \          / \
        //   B   C        A   B
        p := n.Left
        t := p.Right

        n.Left = t
        p.Right = t.Left
        t.Left = p

        *nPtr = t
        n.Left = t.Right
        t.Right = n

        recalculate = append(recalculate, p, t)
      }
    } else {
      // n.Right is too tall.
      if imbalance(n.Right) < 0 {
        // We must do a right-left rotation.
        //     n              n                 t
        //    / \            / \               / \
        //   A   p    ==>   A   t      ==>   n     p
        //      / \            / \          / \   / \
        //     t   D          B   p        A   B C   D
        //    / \                / \
        //   B   C              C   D
        p := n.Right
        t := p.Left

        n.Right = t
        p.Left = t.Right
        t.Right = p

        *nPtr = t
        n.Right = t.Left
        t.Left = n

        recalculate = append(recalculate, p, t)
      } else {
        // We must do an single left rotation.
        //     n               p
        //    / \             / \
        //   A   p    ==>    n   C
        //      / \         / \
        //     B   C       A   B
        // (Where C is taller than B).
        p := n.Right
        
        *nPtr = p
        n.Right = p.Left
        p.Left = n

        recalculate = append(recalculate, p)
      }
    }

		updateSubtreeInfo(n)
    for _, p := range recalculate {
		  updateSubtreeInfo(p)
    }
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

func prettyPrint(n *node) string {
  var buf bytes.Buffer
  indent := make([]byte, 0, 64)
  prettyPrintNode(n, "-", append(indent, []byte("  ")...), &buf)
  return buf.String()
}

func prettyPrintNode(n *node, typ string, indent []byte, dest *bytes.Buffer) {
  if n == nil {
    dest.WriteString(fmt.Sprintf("+-%s ()\n", typ))
    return
  }
  dest.WriteString(fmt.Sprintf("+-%s %s\n", typ, n.Word.Word))
  
  children := make([]*node, 0, 2)
  if n.Left != nil {
    children = append(children, n.Left)
  }
  if n.Right != nil {
    children = append(children, n.Right)
  }
  
  for i, child := range children {
    dest.Write(indent)
    
    var childIndent []byte
    if i < len(children) - 1 {
      childIndent = append(indent, []byte("| ")...)
    } else {
      childIndent = append(indent, []byte("  ")...)
    }
    
    typ := "R"
    if child == n.Left {
      typ = "L"
    }
    prettyPrintNode(child, typ, childIndent, dest)
  }
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

// Serialization helpers.

func serialize(out io.Writer, n *node) error {
	if n == nil {
		return binary.Write(out, byteOrder, int8(0))
	}

	if err := binary.Write(out, byteOrder, int8(1)); err != nil {
		return err
	}

  if err := binary.Write(out, byteOrder, n.Word.Weight); err != nil {
    return err
  }
	if err := binary.Write(out, byteOrder, int64(len(n.Word.Word))); err != nil {
		return err
	}
	if _, err := out.Write([]byte(n.Word.Word)); err != nil {
		return err
	}

	if err := serialize(out, n.Left); err != nil {
		return err
	}
	if err := serialize(out, n.Right); err != nil {
		return err
	}

	return nil
}

func deserialize(in io.Reader) (*node, error) {
	var tag int8
	if err := binary.Read(in, byteOrder, &tag); err != nil {
		return nil, err
	}
	if tag == 0 {
		return nil, nil
	}
	if tag != 1 {
		return nil, fmt.Errorf("Bad tag: %d", tag)
	}

	var weight int64
	if err := binary.Read(in, byteOrder, &weight); err != nil {
		return nil, err
	}
	
	var wordLen int64
	if err := binary.Read(in, byteOrder, &wordLen); err != nil {
		return nil, err
	}

	var word = make([]byte, wordLen)
	if _, err := in.Read(word); err != nil {
		return nil, err
	}

	left, err := deserialize(in)
	if err != nil {
		return nil, err
	}

	right, err := deserialize(in)
	if err != nil {
		return nil, err
	}

	n := &node{left, right, WeightedWord{string(word), weight}, 0, 0, 0}
	updateSubtreeInfo(n)
	return n, nil
}
