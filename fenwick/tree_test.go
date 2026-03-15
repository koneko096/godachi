package fenwick

import (
	"math"
	"testing"

	"github.com/koneko096/godachi/internal"
)

// ── sum value (for sumTree / RangeQuery) ──────────────────────────────────────

type sumVal int

func (v sumVal) Merge(w internal.MergeableValue) internal.MergeableValue {
	return v + w.(sumVal)
}
func (v sumVal) Inverse() internal.MergeableValue    { return -v }
func (v sumVal) Scale(n int) internal.MergeableValue { return v * sumVal(n) }
func (v sumVal) Apply(agg internal.MergeableValue, sz int) internal.MergeableValue {
	return agg.(sumVal) + v*sumVal(sz)
}
func (v sumVal) Compose(newer internal.MergeableValue) internal.MergeableValue {
	return v + newer.(sumVal)
}

// ── max value (for cmpTree / Tree) ───────────────────────────────────────────

type maxVal int

func (v maxVal) Merge(w internal.MergeableValue) internal.MergeableValue {
	if w.(maxVal) > v {
		return w
	}
	return v
}
func (v maxVal) Identity() internal.ComparableValue { return maxVal(math.MinInt64) }
func (v maxVal) Apply(agg internal.MergeableValue, sz int) internal.MergeableValue {
	return v
}
func (v maxVal) Compose(newer internal.MergeableValue) internal.MergeableValue {
	return newer
}

// ── min value (for cmpTree / Tree) ───────────────────────────────────────────

type minVal int

func (v minVal) Merge(w internal.MergeableValue) internal.MergeableValue {
	if w.(minVal) < v {
		return w
	}
	return v
}
func (v minVal) Identity() internal.ComparableValue { return minVal(math.MaxInt64) }
func (v minVal) Apply(agg internal.MergeableValue, sz int) internal.MergeableValue {
	return v
}
func (v minVal) Compose(newer internal.MergeableValue) internal.MergeableValue {
	return newer
}

// ── sum helpers ───────────────────────────────────────────────────────────────

func newSumT(n int) internal.RangeQuery {
	return NewSumTree(n, sumVal(0))
}
func smod(tree internal.RangeQuery, i, j int, v sumVal) {
	tree.Modify(index(i), index(j), v)
}
func sqry(t *testing.T, tree internal.RangeQuery, i, j int) sumVal {
	t.Helper()
	return tree.Query(index(i), index(j)).(sumVal)
}

// ── cmp helpers ───────────────────────────────────────────────────────────────

func newMaxT(n int) internal.Tree {
	return NewCmpTree(n, maxVal(math.MinInt64))
}
func newMinT(n int) internal.Tree {
	return NewCmpTree(n, minVal(math.MaxInt64))
}
func cupdate(tree internal.Tree, i int, v internal.ValueType) {
	tree.Update(index(i), v)
}
func cfind(t *testing.T, tree internal.Tree, i int) internal.ValueType {
	t.Helper()
	return tree.Find(index(i))
}

// ════════════════════════════════════════════════════════════════════════════
// sumTree tests (RangeQuery)
// ════════════════════════════════════════════════════════════════════════════

func TestSum_RMQ(t *testing.T) {
	tree := newSumT(7)
	smod(tree, 0, 4, 5)
	smod(tree, 2, 6, 2)
	// index 3=7, 4=7, 5=2 → 16
	if got := sqry(t, tree, 3, 5); got != 16 {
		t.Errorf("want 16, got %d", got)
	}
}

func TestSum_PointUpdate(t *testing.T) {
	tree := newSumT(8)
	smod(tree, 3, 3, 10)
	if got := sqry(t, tree, 3, 3); got != 10 {
		t.Errorf("point: want 10, got %d", got)
	}
	if got := sqry(t, tree, 2, 2); got != 0 {
		t.Errorf("left neighbour: want 0, got %d", got)
	}
	if got := sqry(t, tree, 4, 4); got != 0 {
		t.Errorf("right neighbour: want 0, got %d", got)
	}
}

func TestSum_FullRange(t *testing.T) {
	const n = 8
	tree := newSumT(n)
	smod(tree, 0, n-1, 3)
	if got := sqry(t, tree, 0, n-1); got != sumVal(3*n) {
		t.Errorf("full range: want %d, got %d", 3*n, got)
	}
}

func TestSum_NonOverlapping(t *testing.T) {
	tree := newSumT(8)
	smod(tree, 0, 2, 4)
	smod(tree, 5, 7, 6)
	if got := sqry(t, tree, 2, 5); got != 10 {
		t.Errorf("cross-boundary: want 10, got %d", got)
	}
}

func TestSum_RepeatedUpdates(t *testing.T) {
	tree := newSumT(8)
	smod(tree, 1, 4, 3)
	smod(tree, 1, 4, 3)
	smod(tree, 1, 4, 3)
	if got := sqry(t, tree, 1, 4); got != 36 {
		t.Errorf("repeated: want 36, got %d", got)
	}
}

func TestSum_PointSumEqualsRange(t *testing.T) {
	const n = 6
	tree := newSumT(n)
	smod(tree, 0, 5, 1)
	smod(tree, 2, 4, 3) // [1,1,4,4,4,1]
	full := sqry(t, tree, 0, n-1)
	var sum sumVal
	for i := 0; i < n; i++ {
		sum += sqry(t, tree, i, i)
	}
	if full != sum {
		t.Errorf("range %d != sum of points %d", full, sum)
	}
}

func TestSum_Stress(t *testing.T) {
	const n = 16
	brute := make([]int, n)
	tree := newSumT(n)
	ops := []struct{ i, j, v int }{
		{0, 15, 1}, {3, 9, 4}, {0, 3, 2},
		{7, 15, 3}, {5, 5, 10}, {0, 15, 1},
		{4, 11, 5}, {2, 13, 3},
	}
	for _, op := range ops {
		for k := op.i; k <= op.j; k++ {
			brute[k] += op.v
		}
		smod(tree, op.i, op.j, sumVal(op.v))
	}
	queries := [][2]int{
		{0, 15}, {0, 0}, {15, 15}, {3, 9},
		{4, 12}, {7, 8}, {0, 7}, {8, 15},
		{5, 5}, {1, 14},
	}
	for _, q := range queries {
		var want int
		for k := q[0]; k <= q[1]; k++ {
			want += brute[k]
		}
		if got := sqry(t, tree, q[0], q[1]); got != sumVal(want) {
			t.Errorf("stress [%d,%d]: want %d, got %d", q[0], q[1], want, got)
		}
	}
}

// ════════════════════════════════════════════════════════════════════════════
// cmpTree tests — max-BIT
// ════════════════════════════════════════════════════════════════════════════

func TestMax_SingleUpdate(t *testing.T) {
	tree := newMaxT(8)
	cupdate(tree, 3, maxVal(42))
	// prefix max [0..3] should be 42
	if got := cfind(t, tree, 3).(maxVal); got != 42 {
		t.Errorf("want 42, got %d", got)
	}
	// prefix max [0..2] should still be -inf (identity)
	if got := cfind(t, tree, 2).(maxVal); got != maxVal(math.MinInt64) {
		t.Errorf("before index: want identity, got %d", got)
	}
}

func TestMax_PrefixMax(t *testing.T) {
	tree := newMaxT(8)
	// insert values at specific indices
	cupdate(tree, 0, maxVal(3))
	cupdate(tree, 2, maxVal(7))
	cupdate(tree, 4, maxVal(5))
	cupdate(tree, 6, maxVal(9))

	// prefix max [0..1] = 3
	if got := cfind(t, tree, 1).(maxVal); got != 3 {
		t.Errorf("[0..1]: want 3, got %d", got)
	}
	// prefix max [0..2] = 7
	if got := cfind(t, tree, 2).(maxVal); got != 7 {
		t.Errorf("[0..2]: want 7, got %d", got)
	}
	// prefix max [0..4] = 7
	if got := cfind(t, tree, 4).(maxVal); got != 7 {
		t.Errorf("[0..4]: want 7, got %d", got)
	}
	// prefix max [0..6] = 9
	if got := cfind(t, tree, 6).(maxVal); got != 9 {
		t.Errorf("[0..6]: want 9, got %d", got)
	}
}

func TestMax_MonotoneUpdate(t *testing.T) {
	tree := newMaxT(8)
	cupdate(tree, 2, maxVal(5))
	// monotone increase: update to larger value
	cupdate(tree, 2, maxVal(10))
	if got := cfind(t, tree, 2).(maxVal); got != 10 {
		t.Errorf("after increase: want 10, got %d", got)
	}
	// update with smaller value: BIT ignores it via Merge
	cupdate(tree, 2, maxVal(3))
	if got := cfind(t, tree, 2).(maxVal); got != 10 {
		t.Errorf("after smaller update: want 10 (unchanged), got %d", got)
	}
}

func TestMax_AllSameValue(t *testing.T) {
	const n = 6
	tree := newMaxT(n)
	for i := 0; i < n; i++ {
		cupdate(tree, i, maxVal(7))
	}
	for i := 0; i < n; i++ {
		if got := cfind(t, tree, i).(maxVal); got != 7 {
			t.Errorf("index %d: want 7, got %d", i, got)
		}
	}
}

func TestMax_Stress(t *testing.T) {
	const n = 16
	// brute[i] = max value inserted at index i
	brute := make([]int, n)
	for i := range brute {
		brute[i] = math.MinInt64
	}
	tree := newMaxT(n)

	updates := []struct{ i, v int }{
		{3, 10}, {7, 4}, {1, 15}, {5, 8},
		{3, 20}, {0, 6}, {9, 13}, {7, 25},
		{2, 11}, {6, 3},
	}
	for _, u := range updates {
		if u.v > brute[u.i] {
			brute[u.i] = u.v
		}
		cupdate(tree, u.i, maxVal(u.v))
	}

	for q := 0; q < n; q++ {
		// expected = max of brute[0..q]
		want := math.MinInt64
		for k := 0; k <= q; k++ {
			if brute[k] > want {
				want = brute[k]
			}
		}
		if got := cfind(t, tree, q).(maxVal); got != maxVal(want) {
			t.Errorf("prefix max [0..%d]: want %d, got %d", q, want, got)
		}
	}
}

// ════════════════════════════════════════════════════════════════════════════
// cmpTree tests — min-BIT
// ════════════════════════════════════════════════════════════════════════════

func TestMin_SingleUpdate(t *testing.T) {
	tree := newMinT(8)
	cupdate(tree, 3, minVal(42))
	if got := cfind(t, tree, 3).(minVal); got != 42 {
		t.Errorf("want 42, got %d", got)
	}
	if got := cfind(t, tree, 2).(minVal); got != minVal(math.MaxInt64) {
		t.Errorf("before index: want identity, got %d", got)
	}
}

func TestMin_PrefixMin(t *testing.T) {
	tree := newMinT(8)
	cupdate(tree, 0, minVal(9))
	cupdate(tree, 2, minVal(3))
	cupdate(tree, 4, minVal(7))
	cupdate(tree, 6, minVal(1))

	// prefix min [0..1] = 9
	if got := cfind(t, tree, 1).(minVal); got != 9 {
		t.Errorf("[0..1]: want 9, got %d", got)
	}
	// prefix min [0..2] = 3
	if got := cfind(t, tree, 2).(minVal); got != 3 {
		t.Errorf("[0..2]: want 3, got %d", got)
	}
	// prefix min [0..4] = 3
	if got := cfind(t, tree, 4).(minVal); got != 3 {
		t.Errorf("[0..4]: want 3, got %d", got)
	}
	// prefix min [0..6] = 1
	if got := cfind(t, tree, 6).(minVal); got != 1 {
		t.Errorf("[0..6]: want 1, got %d", got)
	}
}

func TestMin_MonotoneUpdate(t *testing.T) {
	tree := newMinT(8)
	cupdate(tree, 2, minVal(10))
	// monotone decrease: update to smaller value
	cupdate(tree, 2, minVal(5))
	if got := cfind(t, tree, 2).(minVal); got != 5 {
		t.Errorf("after decrease: want 5, got %d", got)
	}
	// update with larger value: BIT ignores it via Merge
	cupdate(tree, 2, minVal(20))
	if got := cfind(t, tree, 2).(minVal); got != 5 {
		t.Errorf("after larger update: want 5 (unchanged), got %d", got)
	}
}

func TestMin_Stress(t *testing.T) {
	const n = 16
	brute := make([]int, n)
	for i := range brute {
		brute[i] = math.MaxInt64
	}
	tree := newMinT(n)

	updates := []struct{ i, v int }{
		{3, 10}, {7, 4}, {1, 15}, {5, 8},
		{3, 2}, {0, 6}, {9, 13}, {7, 1},
		{2, 11}, {6, 3},
	}
	for _, u := range updates {
		if u.v < brute[u.i] {
			brute[u.i] = u.v
		}
		cupdate(tree, u.i, minVal(u.v))
	}

	for q := 0; q < n; q++ {
		want := math.MaxInt64
		for k := 0; k <= q; k++ {
			if brute[k] < want {
				want = brute[k]
			}
		}
		if got := cfind(t, tree, q).(minVal); got != minVal(want) {
			t.Errorf("prefix min [0..%d]: want %d, got %d", q, want, got)
		}
	}
}
