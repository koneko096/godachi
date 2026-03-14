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

func (t *tree) FindIt(key internal.KeyType) internal.Iterator {
	return t.root.findnode(key)
}

func (t *tree) Empty() bool {
	return t.root == nil
}

func (t *tree) Iterator() internal.Iterator {
	return t.root.minimum()
}

func (t *tree) Size() int {
	return t.size
}

func (t *tree) Clear() {
	t.root = nil
	t.size = 0
}

func (t *tree) Insert(key internal.KeyType, value internal.ValueType) {
	t.root = t.root.insert(&node{
		key:   key,
		value: value,
	}, nil)
	t.size++
}

func (t *tree) Delete(key internal.KeyType) {
	n := t.root.findnode(key)
	if n == nil {
		return
	}

	p := n.parent
	var r *node
	var infix bool
	for _, c := range []*node{n.left, n.right} {
		if c != nil {
			if r != nil {
				infix = true
			}
			r = c
		}
	}

	if infix {
		r = n.successor()
		rp := r.parent
		go func() {
			if n.right == r {
				r.insertLeft(n.left)
			} else {
				rp.insertLeft(r.right)
				r.insertLeft(n.left)
				r.insertRight(rp)
			}
		}()
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
	t.size--
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
