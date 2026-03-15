package segtree

import (
	"testing"

	"github.com/koneko096/godachi/internal"
)

type value int

func (v value) Merge(w internal.MergeableValue) internal.MergeableValue {
	return v + w.(value)
}
func (v value) Apply(agg internal.MergeableValue, sz int) internal.MergeableValue {
	return agg.(value) + v*value(sz)
}
func (v value) Compose(newer internal.MergeableValue) internal.MergeableValue {
	return v + newer.(value)
}

// ── helpers ──────────────────────────────────────────────────────────────────

func mustQuery(t *testing.T, tree internal.RangeQuery, i, j int) value {
	t.Helper()
	return tree.Query(index(i), index(j)).(value)
}

func mustModify(tree internal.RangeQuery, i, j int, v value) {
	tree.Modify(index(i), index(j), v)
}

// ── 1. Single-element point update and query ──────────────────────────────────

func TestPointUpdate(t *testing.T) {
	tree := NewTree(8, value(0))
	mustModify(tree, 3, 3, 10)
	if got := mustQuery(t, tree, 3, 3); got != value(10) {
		t.Errorf("point update: want 10, got %d", got)
	}
	// neighbours must remain 0
	if got := mustQuery(t, tree, 2, 2); got != value(0) {
		t.Errorf("left neighbour: want 0, got %d", got)
	}
	if got := mustQuery(t, tree, 4, 4); got != value(0) {
		t.Errorf("right neighbour: want 0, got %d", got)
	}
}

// ── 2. Full-range update, full-range query ────────────────────────────────────

func TestFullRangeUpdate(t *testing.T) {
	const n = 8
	tree := NewTree(n, value(0))
	mustModify(tree, 0, n-1, 3)
	if got := mustQuery(t, tree, 0, n-1); got != value(3*n) {
		t.Errorf("full range: want %d, got %d", 3*n, got)
	}
}

// ── 3. Non-overlapping updates, query across boundary ────────────────────────

func TestNonOverlappingUpdates(t *testing.T) {
	tree := NewTree(8, value(0))
	mustModify(tree, 0, 2, 4) // [4,4,4,0,0,0,0,0]
	mustModify(tree, 5, 7, 6) // [4,4,4,0,0,6,6,6]
	// query [2,5] = 4+0+0+6 = 10
	if got := mustQuery(t, tree, 2, 5); got != value(10) {
		t.Errorf("cross-boundary: want 10, got %d", got)
	}
}

// ── 4. Repeated updates on same range ────────────────────────────────────────

func TestRepeatedUpdates(t *testing.T) {
	tree := NewTree(8, value(0))
	mustModify(tree, 1, 4, 3)
	mustModify(tree, 1, 4, 3)
	mustModify(tree, 1, 4, 3)
	// each of 4 elements got +9
	if got := mustQuery(t, tree, 1, 4); got != value(36) {
		t.Errorf("repeated updates: want 36, got %d", got)
	}
}

// ── 5. Point queries after range update ──────────────────────────────────────

func TestPointQueriesAfterRangeUpdate(t *testing.T) {
	tree := NewTree(8, value(0))
	mustModify(tree, 0, 7, 2) // all = 2
	mustModify(tree, 3, 5, 8) // indices 3-5 = 10
	cases := []struct {
		idx  int
		want value
	}{
		{0, 2}, {1, 2}, {2, 2},
		{3, 10}, {4, 10}, {5, 10},
		{6, 2}, {7, 2},
	}
	for _, c := range cases {
		if got := mustQuery(t, tree, c.idx, c.idx); got != c.want {
			t.Errorf("index %d: want %d, got %d", c.idx, c.want, got)
		}
	}
}

// ── 6. Sum of point queries equals full-range query ───────────────────────────

func TestPointSumEqualsRangeQuery(t *testing.T) {
	const n = 6
	tree := NewTree(n, value(0))
	mustModify(tree, 0, 5, 1)
	mustModify(tree, 2, 4, 3) // values: [1,1,4,4,4,1]

	full := mustQuery(t, tree, 0, n-1)
	var sum value
	for i := 0; i < n; i++ {
		sum += mustQuery(t, tree, i, i)
	}
	if full != sum {
		t.Errorf("range query %d != sum of point queries %d", full, sum)
	}
}

// ── 7. Zero-value update is a no-op ──────────────────────────────────────────

func TestZeroValueUpdate(t *testing.T) {
	const n = 8
	tree := NewTree(n, value(0))
	mustModify(tree, 0, n-1, 5) // all = 5
	mustModify(tree, 2, 4, 0)   // no-op
	if got := mustQuery(t, tree, 0, n-1); got != value(5*n) {
		t.Errorf("zero update should be no-op: want %d, got %d", 5*n, got)
	}
}

// ── 8. Size-1 tree ────────────────────────────────────────────────────────────

func TestSizeOneTree(t *testing.T) {
	tree := NewTree(1, value(0))
	mustModify(tree, 0, 0, 99)
	mustModify(tree, 0, 0, 1)
	if got := mustQuery(t, tree, 0, 0); got != value(100) {
		t.Errorf("size-1 tree: want 100, got %d", got)
	}
}

// ── 9. Boundary elements after large update ───────────────────────────────────

func TestBoundaryElements(t *testing.T) {
	const n = 10
	tree := NewTree(n, value(0))
	mustModify(tree, 0, n-1, 7)   // all = 7
	mustModify(tree, 0, 0, 3)     // index 0 = 10
	mustModify(tree, n-1, n-1, 3) // index 9 = 10

	if got := mustQuery(t, tree, 0, 0); got != value(10) {
		t.Errorf("left boundary: want 10, got %d", got)
	}
	if got := mustQuery(t, tree, n-1, n-1); got != value(10) {
		t.Errorf("right boundary: want 10, got %d", got)
	}
	if got := mustQuery(t, tree, 1, n-2); got != value(7*8) {
		t.Errorf("middle: want %d, got %d", 7*8, got)
	}
	if got := mustQuery(t, tree, 0, n-1); got != value(10+7*8+10) {
		t.Errorf("full: want %d, got %d", 10+7*8+10, got)
	}
}

// ── 10. Nested overlapping updates ───────────────────────────────────────────

func TestNestedOverlappingUpdates(t *testing.T) {
	// index:    0  1  2  3  4  5  6  7
	// +1 [0,7]: 1  1  1  1  1  1  1  1
	// +2 [1,6]: 1  3  3  3  3  3  3  1
	// +4 [2,5]: 1  3  7  7  7  7  3  1
	// +8 [3,4]: 1  3  7 15 15  7  3  1
	tree := NewTree(8, value(0))
	mustModify(tree, 0, 7, 1)
	mustModify(tree, 1, 6, 2)
	mustModify(tree, 2, 5, 4)
	mustModify(tree, 3, 4, 8)

	expected := []value{1, 3, 7, 15, 15, 7, 3, 1}
	for i, want := range expected {
		if got := mustQuery(t, tree, i, i); got != want {
			t.Errorf("nested index %d: want %d, got %d", i, want, got)
		}
	}
	// full sum: 1+3+7+15+15+7+3+1 = 52
	if got := mustQuery(t, tree, 0, 7); got != value(52) {
		t.Errorf("nested full sum: want 52, got %d", got)
	}
}

// ── 11. Stress: brute-force vs segment tree ───────────────────────────────────

func TestStressBruteForce(t *testing.T) {
	const n = 16
	brute := make([]int, n)
	tree := NewTree(n, value(0))

	ops := []struct{ i, j, v int }{
		{0, 15, 1},
		{3, 9, 4},
		{0, 3, 2},
		{7, 15, 3},
		{5, 5, 10},
		{0, 15, 1},
		{4, 11, 5},
		{2, 13, 3},
	}
	for _, op := range ops {
		for k := op.i; k <= op.j; k++ {
			brute[k] += op.v
		}
		mustModify(tree, op.i, op.j, value(op.v))
	}

	queries := [][2]int{
		{0, 15}, {0, 0}, {15, 15},
		{3, 9}, {4, 12}, {7, 8},
		{0, 7}, {8, 15}, {5, 5},
		{1, 14},
	}
	for _, q := range queries {
		var expected int
		for k := q[0]; k <= q[1]; k++ {
			expected += brute[k]
		}
		if got := mustQuery(t, tree, q[0], q[1]); got != value(expected) {
			t.Errorf("stress [%d,%d]: want %d, got %d", q[0], q[1], expected, got)
		}
	}
}
