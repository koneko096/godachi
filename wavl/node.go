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

// Next returns the node's successor as an iterator
func (n *node) Next() internal.Iterator {
	return n.successor()
}

func (n *node) preorder() {
	fmt.Printf("(%v %v)", n.key, n.value)
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

func (n *node) insert(c *node, p *node) *node {
	if c == nil {
		return n
	}
	if n == nil {
		c.parent = p
		return c
	}
	if c.key.LessThan(n.key) {
		n.left = n.left.insert(c, n)
	} else {
		n.right = n.right.insert(c, n)
	}
	r := n.rebalance()
	return r
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

func (n *node) rebalance() *node {
	if n == nil {
		return n
	}

	nl := n.left.rebalance()
	nr := n.right.rebalance()
	n.refreshRank()

	for {
		if nr.Rank() > nl.Rank()+1 {
			nrr := nr.right
			nrl := nr.left
			if nrl.Rank() > nrr.Rank() {
				nr = nrl.rotateRight(nr)
			} else {
				n = nr.rotateLeft(n)
				nl = n.left
				nr = n.right
			}
		} else if nl.Rank() > nr.Rank()+1 {
			nlr := nl.right
			nll := nl.left
			if nlr.Rank() > nll.Rank() {
				nl = nlr.rotateLeft(nl)
			} else {
				n = nl.rotateRight(n)
				nr = n.right
				nl = n.left
			}
		} else {
			break
		}
		n.refreshRank()
	}
	return n
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

func (n *node) refreshRank() {
	n.rank = max(n.left.Rank(), n.right.Rank()) + 1
}

func (n *node) rotateRight(p *node) *node {
	p.insertLeft(n.right)
	n.parent = p.parent
	n.insertRight(p)
	p.refreshRank()
	n.refreshRank()
	return n
}

func (n *node) rotateLeft(p *node) *node {
	p.insertRight(n.left)
	n.parent = p.parent
	n.insertLeft(p)
	p.refreshRank()
	n.refreshRank()
	return n
}
