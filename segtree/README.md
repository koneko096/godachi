# Segment Tree

A generic segment tree implementation in Go supporting range update and range query in **O(log n)** time, backed by lazy propagation.

The tree operates on any type that satisfies the `internal.MergeableValueType` interface, making it reusable for sum, min, max, XOR, and other associative operations.

```
import "github.com/koneko096/godachi/internal"

var tree internal.RangeQuery = NewTree(n, value(0))
```

---

## Range Query Operations

The tree exposes two operations through the `internal.RangeQuery` interface.

```go
type RangeQuery interface {
    Modify(i, j KeyType, v MergeableValueType)
    Query(i, j KeyType) MergeableValueType
}
```

Both `i` and `j` are **inclusive** bounds. Indices are wrapped in the `index` type which satisfies `internal.KeyType`.

---

### Query

```go
func (t *tree) Query(i, j internal.KeyType) internal.MergeableValueType
```

Returns the merged value of all elements in the closed range `[i, j]`.

**How it works:**

The tree is stored as a 1-indexed array of size `4n`. Each internal node holds the merged value of its entire subtree range. Starting from the leaves layer, the algorithm ascends toward the root, collecting results only for nodes whose range falls within `[l, r]`. Before descending into children, any pending lazy value is pushed down so child nodes reflect all prior updates.

```
Query([3, 5]) on a size-7 tree

              [0,6]=21
             /        \
        [0,3]=14      [4,6]=7
        /     \       /     \
    [0,1]=2 [2,3]=12 [4,5]=5 [6,6]=2
    /   \   /    \   /   \
  [0]=1[1]=1[2]=7[3]=5[4]=2[5]=3

  Collected nodes: [3,3] + [4,5] = 5 + 5 = 10  (example values)
```

**Complexity:** O(log n)

**Example:**

```go
tree := NewTree(7, value(0))
tree.Modify(index(0), index(4), value(5))
tree.Modify(index(2), index(6), value(2))

ans := tree.Query(index(3), index(5)).(value) // returns 16
```

Breakdown: index 3 = 5+2 = 7, index 4 = 5+2 = 7, index 5 = 2 → 7+7+2 = **16**.

---

### Update

```go
func (t *tree) Modify(i, j internal.KeyType, v internal.MergeableValueType)
```

Applies value `v` to every element in the closed range `[i, j]` using `Merge`.

**How it works:**

The algorithm descends recursively. When a node's range is fully covered by `[l, r]`, it applies `v` immediately — merging `v` once per element in that node's range into the node's aggregate value — and records `v` in the node's lazy slot for deferred propagation to children. When a node's range is only partially covered, any existing lazy value is first pushed down to children, then the algorithm recurses into the relevant subtree(s) and recomputes the node's value from its children on the way back up.

```
Modify([2, 6], value(2)) on size-7 tree

              root
             /    \
           [0,3]  [4,6]  ← fully covered: apply +2, store lazy
           /   \
        [0,1] [2,3]      ← [2,3] fully covered: apply +2, store lazy
                         ← [0,1] out of range: skip
```

**Complexity:** O(log n)

**Example:**

```go
tree := NewTree(7, value(0))
tree.Modify(index(0), index(6), value(1)) // add 1 to all
tree.Modify(index(2), index(4), value(3)) // add 3 to indices 2-4
// result: [1, 1, 4, 4, 4, 1, 1]
```

---

## Lazy Propagation

Lazy propagation defers the cost of range updates. Instead of immediately updating every leaf in a range, a pending value is stored on internal nodes and only propagated when a node must be split — that is, when a later query or update needs to distinguish between its two children.

### Node value invariant

Every node stores the **total merged value of all elements in its range**, accounting for all updates applied so far — including ones recorded only in ancestor lazy slots. When a lazy value is pushed down, the children's aggregate values and lazy slots are updated to reflect it, and the parent's lazy slot is cleared.

```go
// Applying lazy to children during pushDown
leftLen  := mid - nodeL + 1
rightLen := nodeR - mid

arr[2*node]   = arr[2*node].Merge(mergeN(lazy[node], leftLen))
lazy[2*node]  = lazy[2*node].Merge(lazy[node])

arr[2*node+1]  = arr[2*node+1].Merge(mergeN(lazy[node], rightLen))
lazy[2*node+1] = lazy[2*node+1].Merge(lazy[node])

lazy[node] = initVal  // clear
```

### mergeN

Because a node's aggregate value represents the sum over its entire range, applying a per-element delta `v` to a node covering `k` elements requires merging `v` exactly `k` times:

```go
func mergeN(v MergeableValueType, n int, init MergeableValueType) MergeableValueType {
    result := init
    for i := 0; i < n; i++ {
        result = result.Merge(v)
    }
    return result
}
```

For addition this is equivalent to `v * k`. For other monoids (min, max, XOR) the behaviour follows from the `Merge` implementation.

### Custom value types

Any type can be used as the value type by implementing `MergeableValueType`:

```go
type MergeableValueType interface {
    Merge(MergeableValueType) MergeableValueType
}
```

| Use case   | `Merge` definition       | `initVal`    |
|------------|--------------------------|--------------|
| Range sum  | `a + b`                  | `0`          |
| Range min  | `min(a, b)`              | `+∞`         |
| Range max  | `max(a, b)`              | `-∞`         |
| Range XOR  | `a ^ b`                  | `0`          |

> **Note:** For min/max trees, `mergeN(v, k)` collapses to just `v` regardless of `k`, because applying the same value to every element in a range leaves the aggregate unchanged beyond the first application. Adjust `mergeN` accordingly for your monoid.

### Complexity summary

| Operation | Time       | Space  |
|-----------|------------|--------|
| Build     | O(n)       | O(n)   |
| Update    | O(log n)   | —      |
| Query     | O(log n)   | —      |
