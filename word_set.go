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
  if self.root == nil {
    return
  }
  if red(self.root) {
    log.Fatal("Root is red")
  }
  check(self.root)
}

func (self WordSet) Add(word WeightedWord) {
  var newNode *node = nil

  if self.root == nil {
    newNode = newLeafNode(nil, word)
    self.root = newNode
  } else {
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
        continue
      }
      
      // c > 0
      if cur.Right == nil {
        newNode = newLeafNode(cur, word)
        cur.Right = newNode
        break
      }
      cur = cur.Right
      continue
    }
  }

  updateSubtreeCounts(newNode)
  
  // We must now rebalance after the insertion of 'newNode'. This logic is
  // based on code from https://en.wikipedia.org/wiki/Red%E2%80%93black_tree.
  insertCase1(newNode)
}

///////////////////////////////////////////////////////////////////////////////
// HELPERS

// Returns the black depth of the subtree rooted at this node.
func check(n *node) int {
  if n == nil {
    // Nil nodes are considered black and so have a black depth of 1.
    return 1
  }
  
  if red(n) && (red(n.Left) || red(n.Right)) {
    log.Fatal("Red node has a red child")
  }
  
  if n.Left != nil {
    if n.Left.Parent != n {
      log.Fatal("Bad parent link")
    }
  }
  if n.Right != nil {
    if n.Right.Parent != n {
      log.Fatal("Bad parent link")
    }
  }
  
  if n.SubtreeNodes != subtreeNodes(n.Left) + subtreeNodes(n.Right) + 1 {
    log.Fatal("Bad SubtreeNodes")
  }
  if n.SubtreeWeight != subtreeWeight(n.Left) + subtreeWeight(n.Right) +
     n.Word.Weight {
    log.Fatal("Bad SubtreeWeight")
  }
  
  leftBlackDepth := check(n.Left)
  rightBlackDepth := check(n.Right)
  if leftBlackDepth != rightBlackDepth {
    log.Fatal("Unequal black depths")
  }
  
  blackDepth := leftBlackDepth
  if black(n) {
    blackDepth++
  }
  return blackDepth
}

// Updates subtree counts at 'n' and all of its ancestors.
func updateSubtreeCounts(n *node) {
  for n != nil {
    n.SubtreeNodes = subtreeNodes(n.Left) + subtreeNodes(n.Right) + 1
    n.SubtreeWeight = subtreeWeight(n.Left) + subtreeWeight(n.Right) +
                      n.Word.Weight
    n = n.Parent
  }
}

func black(n *node) bool {
  if n == nil {
    return true
  }
  return n.Black
}

func red(n *node) bool {
  return !black(n)
}

func subtreeNodes(n *node) int64 {
  if n == nil {
    return 0
  }
  return n.SubtreeNodes
}

func subtreeWeight(n *node) int64 {
  if n == nil {
    return 0
  }
  return n.SubtreeWeight
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
