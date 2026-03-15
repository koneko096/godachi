package treap

import (
	"math/rand/v2"

	"github.com/koneko096/godachi/internal"
)

type accumulator func(x, y internal.ValueType) internal.ValueType
type modifier func(v, x internal.ValueType)

type rangeTree struct {
	root *node
	acc  accumulator
	mod  modifier
}

func NewRangeQueryTree() *rangeTree {
	return &rangeTree{}
}

func (t *rangeTree) Empty() bool {
	return t.root == nil
}

func (t *rangeTree) Size() int {
	return t.root.size()
}

func (t *rangeTree) Clear() {
	t.root = nil
}

func (t *rangeTree) Find(k internal.KeyType) internal.MergeableValue {
	n := t.root.findnode(k)
	if n == nil {
		return nil
	}
	return n.value.(internal.MergeableValue)
}

func (t *rangeTree) Insert(k internal.KeyType, v internal.MergeableValue) {
	t.Delete(k)
	n := &node{
		priority: key(rand.IntN(1 << 30)),
		key:      k,
		value:    v,
		agg:      v,
		sz:       1,
	}
	l, r := t.root.splitRight(k) // l = keys < k, r = keys >= k
	t.root = merge(merge(l, n), r)
}

func (t *rangeTree) Delete(k internal.KeyType) {
	l, mr := t.root.splitRight(k)
	_, r := mr.splitLeft(k)
	t.root = merge(l, r)
}

// Query returns the merged aggregate over the closed range [lo, hi].
// Returns nil if no keys fall within the range.
func (t *rangeTree) Query(lo, hi internal.KeyType) internal.MergeableValue {
	l, mr := t.root.splitRight(lo) // l = keys < lo
	m, r := mr.splitLeft(hi)       // m = keys in [lo, hi], r = keys > hi

	var result internal.MergeableValue
	if m != nil {
		result = m.agg
	}

	t.root = merge(l, merge(m, r))
	return result
}

// Modify applies tag to every node whose key is in [lo, hi].
func (t *rangeTree) Modify(lo, hi internal.KeyType, tag internal.MergeableValue) {
	l, mr := t.root.splitRight(lo)
	m, r := mr.splitLeft(hi)

	m.applyTag(tag)

	t.root = merge(l, merge(m, r))
}
