// Defines a container for words which tracks word weights and allows random
// sampling. Words are stored in a sorted red/black tree, with each internal
// node storing the number and total weight of its descendant leaves.

package dorkalonius

import (
  "log"
  "strings"
)

type WeightedWord struct {
  Word   string
  Weight int64
}

type node struct {
  // Red or black?
  Black  bool
    
  // Parent/child pointers.
  Parent *node
  Left   *node
  Right  *node
    
  // The Weighted represented by this node.
  Word   WeightedWord
  
  // Information about the subtree rooted at this node, which includes
  // this node.
  SubtreeNodes  int64
  SubtreeWeight int64
}

func newLeafNode(parent *node, word WeightedWord) *node {
  // Newly inserted nodes start out as red.
  return &node{false, parent, nil, nil, word, 1, word.Weight}
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
  if self.root != nil {
    if red(self.root) {
      log.Fatal("Root is red")
    }
    self.root.Check()
  }
}

func (self WordSet) Add(word WeightedWord) {
  if self.root == nil {
    self.root = newLeafNode(nil, word)
    self.root.Black = true
    return
  }
  
  var newNode *node = nil
  cur := self.root
  for {
    // No matter what happens, we are adding this weight somewhere.
    cur.SubtreeWeight += word.Weight

    c := strings.Compare(word.Word, cur.Word.Word)
    if c == 0 {
      cur.Word.Weight += word.Weight
      // No need to insert any new nodes, so no need to rebalance. We
      // are done.
      return
    }
    
    // We will have to add a new node.
    cur.SubtreeNodes += 1

    if c < 0 {
      if cur.Left == nil {
        newNode = newLeafNode(cur, word)
        cur.Left = newNode
        break
      }
      cur = cur.Left
    } else {
      if cur.Right == nil {
        newNode = newLeafNode(cur, word)
        cur.Right = newNode
        break
      }
      cur = cur.Right
    }
  }
  
  // We must now rebalance after the insertion of 'newNode'.
}

///////////////////////////////////////////////////////////////////////////////
// HELPERS

// Returns the black depth of the subtree rooted at this node.
func (self *node) Check() int {
  if self.Left == nil && self.Right == nil {
    if self.SubtreeNodes != 1 {
      log.Fatalf("Leaf has wrong SubtreeNodes: %d", self.SubtreeNodes)
    }
    if self.SubtreeWeight != self.Word.Weight {
      log.Fatal("Leaf has wrong SubtreeWeight")
    }
  }
  
  if red(self) && (red(self.Left) || red(self.Right)) {
    log.Fatal("Red node has a red child")
  }
  
  // Nil nodes are considered black and so have a black depth of 1.
  var leftBlackDepth int = 1
  var rightBlackDepth int = 1
  if self.Left != nil {
    leftBlackDepth = self.Left.Check()
  }
  if self.Right != nil {
    rightBlackDepth = self.Right.Check()
  }
  if leftBlackDepth != rightBlackDepth {
    log.Fatal("Unequal black depths")
  }
  
  blackDepth := leftBlackDepth
  if black(self) {
    blackDepth++
  }
  return blackDepth
}

// Red-black tree insertion based on code from
// https://en.wikipedia.org/wiki/Red%E2%80%93black_tree.

func black(n *node) bool {
  if n == nil {
    return true
  }
  return n.Black
}

func red(n *node) bool {
  return !black(n)
}

func grandparent(n *node) *node {
  if n != nil && n.Parent != nil {
    return n.Parent.Parent;
  }
  return nil;
}

func uncle(n *node) *node {
  g := grandparent(n)
  if g == nil {
    return nil
  }
  if n.Parent == g.Left {
    return g.Right;
  }
  return g.Left;
}
