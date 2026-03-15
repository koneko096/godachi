package segtree

import (
	"github.com/koneko096/godachi/internal"
)

type index int

func (n index) LessThan(b internal.KeyType) bool {
	keyB := b.(index)
	return n < keyB
}

func (n index) Equal(b internal.KeyType) bool {
	keyB := b.(index)
	return n == keyB
}

type tree struct {
	sz      int
	arr     []internal.MergeableValue
	lazy    []internal.MergeableValue
	initVal internal.MergeableValue
}

func NewTree(n int, initVal internal.MergeableValue) *tree {
	arr := make([]internal.MergeableValue, 4*n)
	lazy := make([]internal.MergeableValue, 4*n)
	for i := range arr {
		arr[i] = initVal
		lazy[i] = initVal
	}
	return &tree{sz: n, arr: arr, lazy: lazy, initVal: initVal}
}

func (t *tree) Find(key internal.KeyType) internal.ValueType {
	return t.arr[t.generateIndex(key, 0)]
}

func (t *tree) Update(key internal.KeyType, value internal.ValueType) {
	// t.arr[t.generateIndex(key, 0)] = value
}

func (t *tree) generateIndex(idx internal.KeyType, level int) int {
	ival, ok := idx.(index)
	if !ok {
		return -1
	}
	return int(ival)
}

func (t *tree) Modify(i, j internal.KeyType, v internal.MergeableValue) {
	l := t.generateIndex(i, 0)
	r := t.generateIndex(j, 0)
	if l < 0 || r < 0 || l > r {
		return
	}
	t.modifyRange(1, 0, t.sz-1, l, r, v)
}

func (t *tree) Query(i, j internal.KeyType) internal.MergeableValue {
	l := t.generateIndex(i, 0)
	r := t.generateIndex(j, 0)
	if l < 0 || r < 0 || l > r {
		return t.initVal
	}
	return t.queryRange(1, 0, t.sz-1, l, r)
}

func (t *tree) pushDown(node, leftLen, rightLen int) {
	if t.lazy[node] == t.initVal {
		return
	}
	// Apply to left child
	t.arr[2*node] = t.arr[2*node].Merge(
		mergeN(t.lazy[node], leftLen, t.initVal),
	)
	t.lazy[2*node] = t.lazy[2*node].Merge(t.lazy[node])

	// Apply to right child
	t.arr[2*node+1] = t.arr[2*node+1].Merge(
		mergeN(t.lazy[node], rightLen, t.initVal),
	)
	t.lazy[2*node+1] = t.lazy[2*node+1].Merge(t.lazy[node])

	t.lazy[node] = t.initVal
}

// mergeN merges v with itself n times (i.e. v*n for addition)
func mergeN(v internal.MergeableValue, n int, init internal.MergeableValue) internal.MergeableValue {
	result := init
	for _ = range n {
		result = result.Merge(v)
	}
	return result
}

func (t *tree) modifyRange(node, nodeL, nodeR, l, r int, v internal.MergeableValue) {
	if l > nodeR || r < nodeL {
		return
	}
	if l <= nodeL && nodeR <= r {
		// Entire node range is covered
		t.arr[node] = t.arr[node].Merge(mergeN(v, nodeR-nodeL+1, t.initVal))
		t.lazy[node] = t.lazy[node].Merge(v)
		return
	}
	mid := (nodeL + nodeR) / 2
	t.pushDown(node, mid-nodeL+1, nodeR-mid)
	t.modifyRange(2*node, nodeL, mid, l, r, v)
	t.modifyRange(2*node+1, mid+1, nodeR, l, r, v)
	t.arr[node] = t.arr[2*node].Merge(t.arr[2*node+1])
}

func (t *tree) queryRange(node, nodeL, nodeR, l, r int) internal.MergeableValue {
	if l > nodeR || r < nodeL {
		return t.initVal
	}
	if l <= nodeL && nodeR <= r {
		return t.arr[node]
	}
	mid := (nodeL + nodeR) / 2
	t.pushDown(node, mid-nodeL+1, nodeR-mid)
	left := t.queryRange(2*node, nodeL, mid, l, r)
	right := t.queryRange(2*node+1, mid+1, nodeR, l, r)
	return left.Merge(right)
}
