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
	lazy                internal.ValueType
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
	// push_down <- lazy propagation
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
	// push_down <- lazy propagation
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

// TODO: revive for RMQ
// func (t *node) propagate() {
// 	if t.lazy.IsEmpty() {
// 		return
// 	} // nothing to propagate

// 	if !t.left.IsNil() {
// 		t.left.val += t.lazy
// 		t.left.sum += t.lazy * t.left.size() // if tracking range sum
// 		t.left.lazy += t.lazy
// 	}

// 	if !t.right.IsNil() {
// 		t.right.val += t.lazy
// 		t.right.sum += t.lazy * t.right.size()
// 		t.right.lazy += t.lazy
// 	}

// 	t.lazy = 0
// }

func (t *node) update() {
	if t == nil {
		return
	}

	t.sz = 1 + t.left.size() + t.right.size()

	// TODO: revive for RMQ
	// e.g t.sum  = t.val + sum(t.left) + sum(t.right)
	// e.g t.min  = min(t.val, min(t.left), min(t.right))
}
