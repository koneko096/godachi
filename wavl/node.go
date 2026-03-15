package wavl

import (
	"fmt"

	"github.com/koneko096/godachi/internal"
)

type node struct {
	left, right, parent *node
	key                 internal.KeyType
	value               internal.ValueType
	rank                int
}

// rankDiff returns rank(n) - rank(child).
func rankDiff(n, child *node) int {
	return n.Rank() - child.Rank()
}

func (n *node) promote() { n.rank++ }
func (n *node) demote()  { n.rank-- }

// height returns the structural height of the subtree rooted at n.
// Used by fixRotation to decide rotation type — independent of stored rank.
// nil nodes have height -1 (same convention as rank).
func (n *node) height() int {
	if n == nil {
		return -1
	}
	l := n.left.height()
	r := n.right.height()
	if l > r {
		return l + 1
	}
	return r + 1
}

// Next returns the node's successor as an iterator
func (n *node) Next() internal.Iterator {
	s := n.successor()
	if s == nil {
		return nil
	}
	return s
}

func (n *node) preorder() {
	fmt.Printf("(%v %v): ", n.key, n.value)
	if n.left != nil {
		fmt.Printf("%v's left child is ", n.key)
		n.left.preorder()
	}
	if n.right != nil {
		fmt.Printf("%v's right child is ", n.key)
		n.right.preorder()
	}
}

// findnode finds the node by key and return it, if not exists return nil
func (n *node) findnode(key internal.KeyType) *node {
	if n == nil {
		return nil
	}
	if n.key.Equal(key) {
		return n
	}
	if key.LessThan(n.key) {
		return n.left.findnode(key)
	}
	return n.right.findnode(key)
}

func (n *node) insert(c *node, p *node) (*node, *node) {
	if c == nil {
		return n, nil
	}
	if n == nil {
		c.parent = p
		return c, c
	}
	var inserted *node
	if c.key.LessThan(n.key) {
		n.left, inserted = n.left.insert(c, n)
	} else if n.key.LessThan(c.key) {
		n.right, inserted = n.right.insert(c, n)
	} else {
		return n, nil // duplicate — ignore
	}
	return n, inserted
}

// successor returns the successor of the node
func (x *node) successor() *node {
	if x.right != nil {
		return x.right.minimum()
	}
	y := x.parent
	for y != nil && x == y.right {
		x = y
		y = x.parent
	}
	return y
}

// minimum finds the minimum node of subtree n.
func (n *node) minimum() *node {
	if n == nil {
		return nil
	}
	for n.left != nil {
		n = n.left
	}
	return n
}

// maximum finds the maximum node of subtree n.
func (n *node) maximum() *node {
	if n == nil {
		return nil
	}
	for n.right != nil {
		n = n.right
	}
	return n
}

// treeRoot walks up from n to the root.
func (n *node) treeRoot() *node {
	if n == nil {
		return nil
	}
	for n.parent != nil {
		n = n.parent
	}
	return n
}

func (n *node) Rank() int {
	if n == nil {
		return -1
	}
	return n.rank
}

func (n *node) IsNil() bool {
	return n == nil
}

func (n *node) Key() internal.KeyType {
	return n.key
}

func (n *node) Value() internal.ValueType {
	return n.value
}

func (n *node) Set(v internal.ValueType) {
	n.value = v
}

func (n *node) insertLeft(m *node) {
	n.left = m
	if m == nil {
		return
	}
	m.parent = n
}

func (n *node) insertRight(m *node) {
	n.right = m
	if m == nil {
		return
	}
	m.parent = n
}
