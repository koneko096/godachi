package fenwick

import "github.com/koneko096/godachi/internal"

// ── index key ─────────────────────────────────────────────────────────────────

type index int

func (i index) LessThan(j internal.KeyType) bool {
	jv, ok := j.(index)
	if !ok {
		return false
	}
	return i < jv
}

func (n index) Equal(b internal.KeyType) bool {
	keyB := b.(index)
	return n == keyB
}

// toInternal converts a 0-based external index to a 1-based internal index.
func toInternal(k internal.KeyType) int {
	return int(k.(index)) + 1
}

// ── sum BIT (internal) ────────────────────────────────────────────────────────

type sumBIT struct {
	arr  []internal.MergeableValue
	zero internal.MergeableValue
}

func newSumBIT(n int, zero internal.MergeableValue) *sumBIT {
	arr := make([]internal.MergeableValue, n+2)
	for i := range arr {
		arr[i] = zero
	}
	return &sumBIT{arr: arr, zero: zero}
}

func (b *sumBIT) update(i int, v internal.MergeableValue) {
	n := len(b.arr) - 1
	for ; i <= n; i += i & -i {
		b.arr[i] = b.arr[i].Merge(v)
	}
}

func (b *sumBIT) query(i int) internal.MergeableValue {
	res := b.zero
	for ; i > 0; i -= i & -i {
		res = res.Merge(b.arr[i])
	}
	return res
}

// ── comparable BIT (internal) ─────────────────────────────────────────────────

type cmpBIT struct {
	arr  []internal.ComparableValue
	zero internal.ComparableValue // identity element (−∞ for max, +∞ for min)
}

func newCmpBIT(n int, zero internal.ComparableValue) *cmpBIT {
	arr := make([]internal.ComparableValue, n+1)
	for i := range arr {
		arr[i] = zero
	}
	return &cmpBIT{arr: arr, zero: zero}
}

// update sets position i to Merge(current, v). i is 1-based.
// For correctness, v must be >= current value (max-BIT) or <= (min-BIT).
func (b *cmpBIT) update(i int, v internal.ComparableValue) {
	n := len(b.arr) - 1
	for ; i <= n; i += i & -i {
		b.arr[i] = b.arr[i].Merge(v).(internal.ComparableValue)
	}
}

// query returns the prefix Merge of [1..i]. i is 1-based.
func (b *cmpBIT) query(i int) internal.ComparableValue {
	res := b.zero
	for ; i > 0; i -= i & -i {
		res = res.Merge(b.arr[i]).(internal.ComparableValue)
	}
	return res
}

// ── sumTree: RangeQuery (range update + range query, sum only) ────────────────

// sumTree uses two BITs to support range update and range query in O(log n).
//
// After a sequence of range updates, the prefix sum up to position x is:
//
//	prefix(x) = b1.query(x) * x − b2.query(x)
//
// A range update [l, r] += v is decomposed as:
//
//	b1: +v at l,       −v at r+1
//	b2: +v*(l−1) at l, −v*r at r+1
type sumTree struct {
	n    int
	b1   *sumBIT
	b2   *sumBIT
	zero internal.InvertibleValue
}

// SumTree combines internal.RangeQuery and internal.Tree into one interface,
// since sumTree can answer both range and point operations.
type SumTree interface {
	internal.RangeQuery
	internal.Tree
}

// NewSumTree returns a Fenwick tree of size n for range-update + range-query
// as well as point update and point query.
// zero must implement internal.InvertibleValue and act as the additive identity.
func NewSumTree(n int, zero internal.InvertibleValue) SumTree {
	return &sumTree{
		n:    n,
		b1:   newSumBIT(n, zero),
		b2:   newSumBIT(n, zero),
		zero: zero,
	}
}

func (t *sumTree) prefix(i int) internal.MergeableValue {
	if i <= 0 {
		return t.zero
	}
	scaled := t.b1.query(i).(internal.InvertibleValue).Scale(i)
	offset := t.b2.query(i).(internal.InvertibleValue).Inverse()
	return scaled.Merge(offset)
}

func (t *sumTree) Modify(i, j internal.KeyType, v internal.MergeableValue) {
	l := toInternal(i)
	r := toInternal(j)
	iv := v.(internal.InvertibleValue)

	t.b1.update(l, v)
	t.b1.update(r+1, iv.Inverse())
	t.b2.update(l, iv.Scale(l-1))
	t.b2.update(r+1, iv.Scale(r).(internal.InvertibleValue).Inverse())
}

func (t *sumTree) Query(i, j internal.KeyType) internal.MergeableValue {
	l := toInternal(i)
	r := toInternal(j)
	right := t.prefix(r)
	left := t.prefix(l - 1).(internal.InvertibleValue).Inverse()
	return right.Merge(left)
}

// Update is a point modify: adds value at a single index.
// Delegates to Modify(key, key, value).
func (t *sumTree) Update(key internal.KeyType, value internal.ValueType) {
	t.Modify(key, key, value.(internal.MergeableValue))
}

// Find is a point query: returns the aggregate value at a single index.
// Delegates to Query(key, key).
func (t *sumTree) Find(key internal.KeyType) internal.ValueType {
	return t.Query(key, key)
}

// ── cmpTree: Tree (point update + prefix query, max/min) ─────────────────────

// cmpTree wraps a single comparable BIT and implements internal.Tree.
// It supports:
//   - Update(i, v): sets index i to Merge(current[i], v) — monotone only
//   - Find(i):      prefix Merge of [0..i] — i.e. max/min of all indices <= i
//
// Range query [l, r] is NOT supported because max/min has no inverse.
type cmpTree struct {
	n    int
	b    *cmpBIT
	zero internal.ComparableValue
}

// NewCmpTree returns a Fenwick tree of size n for point update + prefix query.
// zero must implement ComparableValue and be the identity for Merge
// (e.g. math.MinInt for max-BIT, math.MaxInt for min-BIT).
func NewCmpTree(n int, zero internal.ComparableValue) internal.Tree {
	return &cmpTree{
		n:    n,
		b:    newCmpBIT(n, zero),
		zero: zero,
	}
}

// Update applies Merge(current, value) at key.
// For max-BIT: only meaningful when value > current (monotone increase).
// For min-BIT: only meaningful when value < current (monotone decrease).
func (t *cmpTree) Update(key internal.KeyType, value internal.ValueType) {
	i := toInternal(key)
	t.b.update(i, value.(internal.ComparableValue))
}

// Find returns the prefix Merge of all elements in [0..key].
// For max-BIT: the maximum value among indices [0..key].
// For min-BIT: the minimum value among indices [0..key].
func (t *cmpTree) Find(key internal.KeyType) internal.ValueType {
	i := toInternal(key)
	return t.b.query(i)
}
