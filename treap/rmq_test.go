package treap

import (
	"testing"

	"github.com/koneko096/godachi/internal"
)

// ─────────────────────────────────────────────
//  Value types
// ─────────────────────────────────────────────

// sumVal: range-add / sum-query
//
//	Merge(a,b)      = a+b
//	Apply(agg,sz)   = agg + tag*sz   (add tag to every element)
//	Compose(newer)  = self + newer   (additive tags stack)
type sumVal int

func (v sumVal) Merge(o internal.MergeableValue) internal.MergeableValue { return v + o.(sumVal) }
func (v sumVal) Apply(agg internal.MergeableValue, sz int) internal.MergeableValue {
	return agg.(sumVal) + v*sumVal(sz)
}
func (v sumVal) Compose(newer internal.MergeableValue) internal.MergeableValue {
	return v + newer.(sumVal)
}

// assignMin: range-assign / min-query
//
//	Merge(a,b)      = min(a,b)
//	Apply(agg,_)    = self           (overwrite regardless of old value)
//	Compose(newer)  = newer          (later assignment wins)
type assignMin int

func (v assignMin) Merge(o internal.MergeableValue) internal.MergeableValue {
	if o.(assignMin) < v {
		return o
	}
	return v
}
func (v assignMin) Apply(_ internal.MergeableValue, _ int) internal.MergeableValue { return v }
func (v assignMin) Compose(newer internal.MergeableValue) internal.MergeableValue  { return newer }

// ─────────────────────────────────────────────
//  Helpers
// ─────────────────────────────────────────────

func newSumTree() *rangeTree { return NewRangeQueryTree() }
func newMinTree() *rangeTree { return NewRangeQueryTree() }

// insert keyed 1..n with values 1..n into a sum tree.
func fillSum(t *rangeTree, n int) {
	for i := 1; i <= n; i++ {
		t.Insert(key(i), sumVal(i))
	}
}

func assertSum(t *testing.T, tr *rangeTree, lo, hi, want int) {
	t.Helper()
	got := tr.Query(key(lo), key(hi))
	if got == nil || got.(sumVal) != sumVal(want) {
		t.Fatalf("Query(%d,%d): want %d, got %v", lo, hi, want, got)
	}
}

func assertMin(t *testing.T, tr *rangeTree, lo, hi, want int) {
	t.Helper()
	got := tr.Query(key(lo), key(hi))
	if got == nil || got.(assignMin) != assignMin(want) {
		t.Fatalf("Query(%d,%d): want %d, got %v", lo, hi, want, got)
	}
}

// ─────────────────────────────────────────────
//  Empty / single-node edge cases
// ─────────────────────────────────────────────

func TestEmpty(t *testing.T) {
	tr := newSumTree()
	if !tr.Empty() {
		t.Fatal("new tree should be empty")
	}
	if tr.Size() != 0 {
		t.Fatal("new tree size should be 0")
	}
	if got := tr.Query(key(1), key(5)); got != nil {
		t.Fatalf("Query on empty tree: want nil, got %v", got)
	}
}

func TestSingleNode(t *testing.T) {
	tr := newSumTree()
	tr.Insert(key(3), sumVal(42))

	if tr.Size() != 1 {
		t.Fatalf("size: want 1, got %d", tr.Size())
	}
	assertSum(t, tr, 3, 3, 42)

	// range containing only that key
	assertSum(t, tr, 1, 5, 42)

	// range not containing that key
	if got := tr.Query(key(4), key(9)); got != nil {
		t.Fatalf("miss query: want nil, got %v", got)
	}
}

// ─────────────────────────────────────────────
//  Query boundary correctness
// ─────────────────────────────────────────────

// Verifies that lo and hi are both inclusive.
func TestQueryBoundariesInclusive(t *testing.T) {
	tr := newSumTree()
	fillSum(tr, 5) // keys 1..5, values 1..5

	// exact endpoints must be included
	assertSum(t, tr, 1, 1, 1)
	assertSum(t, tr, 5, 5, 5)
	assertSum(t, tr, 1, 5, 15)

	// one inside each boundary
	assertSum(t, tr, 2, 4, 9)

	// query just outside populated range returns nil
	if got := tr.Query(key(6), key(9)); got != nil {
		t.Fatalf("out-of-range query: want nil, got %v", got)
	}
}

// ─────────────────────────────────────────────
//  Insert / Delete
// ─────────────────────────────────────────────

func TestInsertReplacesExistingKey(t *testing.T) {
	tr := newSumTree()
	tr.Insert(key(2), sumVal(10))
	tr.Insert(key(2), sumVal(99)) // replace

	if tr.Size() != 1 {
		t.Fatalf("size after replace: want 1, got %d", tr.Size())
	}
	assertSum(t, tr, 2, 2, 99)
}

func TestDeleteMiddle(t *testing.T) {
	tr := newSumTree()
	fillSum(tr, 5) // [1,2,3,4,5]

	tr.Delete(key(3))

	if tr.Size() != 4 {
		t.Fatalf("size after delete: want 4, got %d", tr.Size())
	}
	assertSum(t, tr, 1, 5, 12) // 1+2+4+5
	assertSum(t, tr, 1, 2, 3)
	assertSum(t, tr, 4, 5, 9)
}

func TestDeleteEndpoints(t *testing.T) {
	tr := newSumTree()
	fillSum(tr, 5)

	tr.Delete(key(1))
	assertSum(t, tr, 1, 5, 14) // 2+3+4+5

	tr.Delete(key(5))
	assertSum(t, tr, 1, 5, 9) // 2+3+4
}

func TestDeleteAbsent(t *testing.T) {
	tr := newSumTree()
	fillSum(tr, 3)
	tr.Delete(key(99)) // no-op
	assertSum(t, tr, 1, 3, 6)
}

// ─────────────────────────────────────────────
//  Find
// ─────────────────────────────────────────────

func TestFind(t *testing.T) {
	tr := newSumTree()
	fillSum(tr, 5)

	for i := 1; i <= 5; i++ {
		got := tr.Find(key(i))
		if got == nil || got.(sumVal) != sumVal(i) {
			t.Fatalf("Find(%d): want %d, got %v", i, i, got)
		}
	}
	if got := tr.Find(key(99)); got != nil {
		t.Fatalf("Find(absent): want nil, got %v", got)
	}
}

// Find must reflect a prior Modify (lazy tag must be flushed on descent).
func TestFindAfterModify(t *testing.T) {
	tr := newSumTree()
	fillSum(tr, 5)

	tr.Modify(key(1), key(5), sumVal(10)) // all += 10

	for i := 1; i <= 5; i++ {
		want := sumVal(i + 10)
		got := tr.Find(key(i))
		if got == nil || got.(sumVal) != want {
			t.Fatalf("Find(%d) after Modify: want %d, got %v", i, want, got)
		}
	}
}

// ─────────────────────────────────────────────
//  Modify (range-add / sum)
// ─────────────────────────────────────────────

func TestModifyFullRange(t *testing.T) {
	tr := newSumTree()
	fillSum(tr, 5) // [1,2,3,4,5]  sum=15

	tr.Modify(key(1), key(5), sumVal(10)) // all += 10 → [11,12,13,14,15]
	assertSum(t, tr, 1, 5, 65)
}

func TestModifySubRange(t *testing.T) {
	tr := newSumTree()
	fillSum(tr, 5) // [1,2,3,4,5]

	tr.Modify(key(2), key(4), sumVal(10)) // [1,12,13,14,5]
	assertSum(t, tr, 1, 5, 45)
	assertSum(t, tr, 2, 4, 39)
	assertSum(t, tr, 1, 1, 1) // untouched
	assertSum(t, tr, 5, 5, 5) // untouched
}

func TestModifyBoundaryNodes(t *testing.T) {
	tr := newSumTree()
	fillSum(tr, 5)

	// modify only the boundary nodes themselves
	tr.Modify(key(1), key(1), sumVal(100)) // [101,2,3,4,5]
	tr.Modify(key(5), key(5), sumVal(100)) // [101,2,3,4,105]
	assertSum(t, tr, 1, 5, 215)
	assertSum(t, tr, 2, 4, 9) // middle untouched
}

func TestModifyNestedOverlapping(t *testing.T) {
	tr := newSumTree()
	fillSum(tr, 5) // [1,2,3,4,5]

	tr.Modify(key(1), key(5), sumVal(10))  // [11,12,13,14,15]
	tr.Modify(key(2), key(4), sumVal(100)) // [11,112,113,114,15]
	assertSum(t, tr, 1, 5, 365)
	assertSum(t, tr, 2, 4, 339)
	assertSum(t, tr, 1, 1, 11)
	assertSum(t, tr, 5, 5, 15)
}

// Modify followed by Query followed by another Modify must stay consistent.
func TestModifyQueryModify(t *testing.T) {
	tr := newSumTree()
	fillSum(tr, 4) // [1,2,3,4]

	tr.Modify(key(1), key(4), sumVal(1)) // [2,3,4,5]  sum=14
	assertSum(t, tr, 1, 4, 14)

	tr.Modify(key(2), key(3), sumVal(10)) // [2,13,14,5]  sum=34
	assertSum(t, tr, 1, 4, 34)
	assertSum(t, tr, 1, 2, 15)
	assertSum(t, tr, 3, 4, 19)
}

// ─────────────────────────────────────────────
//  Modify (range-assign / min)
// ─────────────────────────────────────────────

func TestAssignMinBasic(t *testing.T) {
	tr := newMinTree()
	for i := 1; i <= 5; i++ {
		tr.Insert(key(i), assignMin(i*10)) // [10,20,30,40,50]
	}

	assertMin(t, tr, 1, 5, 10)
	assertMin(t, tr, 3, 5, 30)

	tr.Modify(key(2), key(4), assignMin(7)) // [10,7,7,7,50]
	assertMin(t, tr, 1, 5, 7)
	assertMin(t, tr, 1, 1, 10)
	assertMin(t, tr, 5, 5, 50)
}

func TestAssignMinOverwrite(t *testing.T) {
	tr := newMinTree()
	for i := 1; i <= 5; i++ {
		tr.Insert(key(i), assignMin(i*10))
	}

	tr.Modify(key(1), key(5), assignMin(3)) // all become 3
	assertMin(t, tr, 1, 5, 3)

	// overwrite again with a higher value — min query still returns 3
	// because assign replaces, not adds
	tr.Modify(key(3), key(3), assignMin(99))
	assertMin(t, tr, 1, 5, 3)
	assertMin(t, tr, 3, 3, 99)
}

// Assign then assign again (compose): later tag must win.
func TestAssignComposeLatestWins(t *testing.T) {
	tr := newMinTree()
	for i := 1; i <= 3; i++ {
		tr.Insert(key(i), assignMin(100))
	}

	// two Modifys on the same range without an intervening Query
	// so the second tag is composed onto the first lazy tag
	tr.Modify(key(1), key(3), assignMin(50))
	tr.Modify(key(1), key(3), assignMin(20))
	assertMin(t, tr, 1, 3, 20)

	got := tr.Find(key(2))
	if got == nil || got.(assignMin) != 20 {
		t.Fatalf("Find after composed assign: want 20, got %v", got)
	}
}

// ─────────────────────────────────────────────
//  Size / Clear
// ─────────────────────────────────────────────

func TestSizeAndClear(t *testing.T) {
	tr := newSumTree()
	fillSum(tr, 5)

	if tr.Size() != 5 {
		t.Fatalf("size: want 5, got %d", tr.Size())
	}

	tr.Clear()

	if !tr.Empty() {
		t.Fatal("after Clear, tree should be empty")
	}
	if tr.Size() != 0 {
		t.Fatalf("after Clear, size: want 0, got %d", tr.Size())
	}
	if got := tr.Query(key(1), key(5)); got != nil {
		t.Fatalf("Query after Clear: want nil, got %v", got)
	}
}

// ─────────────────────────────────────────────
//  Aggregate consistency after structural changes
// ─────────────────────────────────────────────

// agg must remain correct after interleaved Insert, Delete, Modify, Query.
func TestAggConsistencyInterleavedOps(t *testing.T) {
	tr := newSumTree()
	fillSum(tr, 6) // [1,2,3,4,5,6]  sum=21

	tr.Modify(key(2), key(5), sumVal(10)) // [1,12,13,14,15,6]  sum=61
	assertSum(t, tr, 1, 6, 61)

	tr.Delete(key(3)) // [1,12,14,15,6]  sum=48
	assertSum(t, tr, 1, 6, 48)

	tr.Insert(key(3), sumVal(100)) // [1,12,100,14,15,6]  sum=148
	assertSum(t, tr, 1, 6, 148)
	assertSum(t, tr, 2, 4, 126) // 12+100+14
}

// Repeatedly query the same range to confirm the tree is restored correctly.
func TestQueryIsNonDestructive(t *testing.T) {
	tr := newSumTree()
	fillSum(tr, 5)

	for i := 0; i < 5; i++ {
		assertSum(t, tr, 1, 5, 15)
		assertSum(t, tr, 2, 4, 9)
	}
}
