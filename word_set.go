// Defines a container for words which tracks word weights and allows random
// sampling. Words are stored in a sorted AVL tree, with each internal
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
  // Parent/child pointers.
  Parent *node // TODO: can we drop this?
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
  return &node{parent, nil, nil, word, 1, word.Weight}
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
  
  // We must now rebalance after the insertion of 'newNode'. TODO:
}

///////////////////////////////////////////////////////////////////////////////
// HELPERS

// Returns the height of the subtree rooted at 'n'.
func check(n *node) int {
  if n == nil {
    return 0
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
  
  leftHeight := check(n.Left)
  rightHeight := check(n.Right)
  imbalance := leftHeight - rightHeight
  if imbalance < -1 || imbalance > 1 {
    log.Fatal("Too much imbalance")
  }
  
  if leftHeight > rightHeight {
    return leftHeight + 1
  }
  return rightHeight + 1
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
