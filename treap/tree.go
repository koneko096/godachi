package treap

import (
	"fmt"
	"math/rand/v2"

	"github.com/koneko096/godachi/internal"
)

type accumulator func(x, y internal.ValueType) internal.ValueType
type modifier func(v, x internal.ValueType)

type tree struct {
	root *node
	acc  accumulator
	mod  modifier
}

func NewTree() *tree {
	return &tree{}
}

func NewRangeTree(acc accumulator, mod modifier) *tree {
	return &tree{acc: acc, mod: mod}
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
	return t.root.size()
}

func (t *tree) Clear() {
	t.root = nil
}

func (t *tree) Insert(k internal.KeyType, value internal.ValueType) {
	p := rand.IntN(10000)
	node := &node{
		priority: key(p),
		key:      k,
		value:    value,
		acc:      t.acc,
		mod:      t.mod,
		sz:       1,
	}
	left, right := t.root.splitLeft(k) // keys <= k | keys > k
	t.root = merge(merge(left, node), right)
}

func (t *tree) Delete(key internal.KeyType) {
	L, midR := t.root.splitRight(key) // L = keys < k,  midR = keys >= k
	_, R := midR.splitLeft(key)       // mid = keys == k, R = keys > k
	// discard mid
	t.root = merge(L, R)
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
