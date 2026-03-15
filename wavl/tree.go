package wavl

import (
	"fmt"

	"github.com/koneko096/godachi/internal"
)

type tree struct {
	root *node
	size int
}

func NewTree() *tree {
	return &tree{}
}

func (t *tree) Find(key internal.KeyType) internal.ValueType {
	n := t.root.findnode(key)
	if n != nil {
		return n.value
	}
	return nil
}

func (t *tree) Update(key internal.KeyType, value internal.ValueType) {
	n := t.root.findnode(key)
	if n != nil {
		n.value = value
	}
}

func (t *tree) FindIt(key internal.KeyType) internal.Iterator {
	n := t.root.findnode(key)
	if n == nil {
		return nil
	}
	return n
}

func (t *tree) Empty() bool {
	return t.root == nil
}

func (t *tree) Iterator() internal.Iterator {
	n := t.root.minimum()
	if n == nil {
		return nil
	}
	return n
}

func (t *tree) Size() int {
	return t.size
}

func (t *tree) Clear() {
	t.root = nil
	t.size = 0
}

func (t *tree) Insert(key internal.KeyType, value internal.ValueType) {
	var n *node
	t.root, n = t.root.insert(&node{
		key:   key,
		value: value,
		rank:  0,
	}, nil)
	if n == nil {
		return // duplicate, nothing to rebalance
	}
	t.size++
	if n.parent == nil {
		// n is the root — no rebalancing needed
		t.root = n
		return
	}
	t.root = rebalanceInsertUp(n.parent)
}

func (t *tree) Delete(key internal.KeyType) {
	n := t.root.findnode(key)
	if n == nil {
		return
	}

	// two children: copy successor payload, redirect deletion to successor
	if n.left != nil && n.right != nil {
		s := n.successor()
		n.key = s.key
		n.value = s.value
		n = s // s has at most one child (right only)
	}

	// 0 or 1 child — your original logic, which is correct
	p := n.parent
	var r *node
	for _, c := range []*node{n.left, n.right} {
		if c != nil {
			r = c
		}
	}
	if r != nil {
		r.parent = p
	}
	if p != nil {
		if p.right == n {
			p.right = r
		} else {
			p.left = r
		}
	} else {
		t.root = r
	}

	n.left, n.right, n.parent = nil, nil, nil
	t.size--

	if p == nil {
		return
	}

	t.root = rebalanceDeleteUp(p)
}

func (t *tree) Preorder() {
	fmt.Println("preorder begin!")
	if t.root != nil {
		t.root.preorder()
	}
	fmt.Println("preorder end!")
}

func (t *tree) Min() internal.ValueType {
	return t.root.minimum().Value()
}

func (t *tree) Max() internal.ValueType {
	return t.root.maximum().Value()
}

// transplant transplants the subtree u and v
func (t *tree) transplant(u, v *node) {}
