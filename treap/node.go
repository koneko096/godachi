package treap

import (
	"fmt"

	"github.com/koneko096/godachi/internal"
)

type node struct {
	left, right, parent *node
	priority            internal.KeyType
	key                 internal.KeyType
	value               internal.ValueType
	agg                 internal.MergeableValue // cached subtree aggregate; nil in plain trees
	lazy                internal.MergeableValue // pending lazy tag; nil = no pending update
	acc                 accumulator
	mod                 modifier
	sz                  int
}

// Next returns the node's successor as an iterator
func (n *node) Next() internal.Iterator {
	return n.successor()
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

func (n *node) preorder() {
	if n == nil {
		return
	}
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

func (n *node) size() int {
	if n == nil {
		return 0
	}
	return n.sz
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

// findnode finds the node by key and return it, if not exists return nil
func (n *node) findnode(key internal.KeyType) *node {
	if n == nil {
		return nil
	}
	n.propagate() // flush lazy so value is current before reading or descending
	if n.key.Equal(key) {
		return n
	}
	if key.LessThan(n.key) {
		return n.left.findnode(key)
	}
	return n.right.findnode(key)
}

// Result: L = keys <= k,  R = keys > k
func (t *node) splitLeft(key internal.KeyType) (*node, *node) {
	if t == nil {
		return nil, nil
	}
	t.propagate()
	if key.LessThan(t.key) {
		// t.key > key → t goes to R, recurse left
		ll, lr := t.left.splitLeft(key)
		t.left = lr
		t.update()
		return ll, t
	} else {
		// t.key <= key → t goes to L, recurse right
		rl, rr := t.right.splitLeft(key)
		t.right = rl
		t.update()
		return t, rr
	}
}

// Result: L = keys < k,  R = keys >= k
func (t *node) splitRight(key internal.KeyType) (*node, *node) {
	if t == nil {
		return nil, nil
	}
	t.propagate()
	if key.LessThan(t.key) || key == t.key { // t.key >= key → t goes to R
		ll, lr := t.left.splitRight(key)
		t.left = lr
		t.update()
		return ll, t
	} else {
		// t.key < key → t goes to L
		rl, rr := t.right.splitRight(key)
		t.right = rl
		t.update()
		return t, rr
	}
}

func merge(L, R *node) *node {
	if L == nil {
		return R
	}
	if R == nil {
		return L
	}
	L.propagate()
	R.propagate()
	if R.priority.LessThan(L.priority) {
		L.right = merge(L.right, R)
		L.update()
		return L
	} else {
		R.left = merge(L, R.left)
		R.update()
		return R
	}
}

// propagate flushes t's lazy tag to its children without descending further.
// Must be called at the top of splitLeft, splitRight, and merge before any
// structural change.  No-op when lazy is nil (plain tree nodes).
func (t *node) propagate() {
	if t == nil || t.lazy == nil {
		return
	}
	t.left.applyTag(t.lazy)
	t.right.applyTag(t.lazy)
	t.lazy = nil
}

// applyTag stamps tag onto t in O(1) without visiting t's children.
func (t *node) applyTag(tag internal.MergeableValue) {
	if t == nil || tag == nil {
		return
	}
	if t.agg != nil {
		t.agg = tag.Apply(t.agg, t.sz)
	}
	if mv, ok := t.value.(internal.MergeableValue); ok {
		t.value = tag.Apply(mv, 1)
	}
	if t.lazy == nil {
		t.lazy = tag
	} else {
		t.lazy = t.lazy.Compose(tag)
	}
}

func (t *node) update() {
	if t == nil {
		return
	}

	t.sz = 1 + t.left.size() + t.right.size()

	// Rebuild agg for range-query trees.  No-op for plain trees where
	// value is not a MergeableValue.
	//
	// Safe to recompute unconditionally: splitLeft, splitRight, and merge
	// all call propagate() before descending into a node, which flushes
	// any pending lazy tag to the children and clears it.  By the time
	// update() is called on the way back up, t.lazy is nil, t.value is
	// current, and each child's agg already reflects its own pending tag.
	if _, ok := t.value.(internal.MergeableValue); !ok {
		return
	}
	t.agg = t.value.(internal.MergeableValue)
	if t.left != nil && t.left.agg != nil {
		t.agg = t.left.agg.Merge(t.agg)
	}
	if t.right != nil && t.right.agg != nil {
		t.agg = t.agg.Merge(t.right.agg)
	}
}
