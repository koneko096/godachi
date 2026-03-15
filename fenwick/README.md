# Fenwick Tree (Binary Indexed Tree)

A generic Fenwick tree (BIT) implementation in Go supporting two distinct operation modes, each backed by a different internal structure and value interface.

```go
import "github.com/koneko096/godachi/fenwick"

// For sum-like types (range update + range query + point update + point query)
var tree fenwick.SumTree = fenwick.NewSumTree(n, value(0))

// For max/min types (monotone point update + prefix query)
var tree internal.Tree = fenwick.NewCmpTree(n, maxValue(math.MinInt64))
```

---

## Constructors

### `NewSumTree`

```go
func NewSumTree(n int, zero InvertibleValue) SumTree
```

Returns a `SumTree` — a composite interface satisfying both `internal.RangeQuery` and `internal.Tree`. Suitable for types where addition, inverse, and scalar multiply are defined (integers, rationals, vectors, etc.).

### `NewCmpTree`

```go
func NewCmpTree(n int, zero ComparableValue) internal.Tree
```

Returns an `internal.Tree` for max or min operations. Supports only monotone point updates and prefix queries. Does not satisfy `internal.RangeQuery`.

---

## Value Interfaces

Two interfaces extend the base `internal` types to unlock the operations each tree mode requires.

### `InvertibleValue` — required by `NewSumTree`

```go
type InvertibleValue interface {
    internal.MergeableValueType
    Inverse() internal.MergeableValueType  // additive inverse:  −v
    Scale(n int) internal.MergeableValueType // repeated merge:   v * n
}
```

**Why `Inverse`:** Range query is derived from two prefix queries:

```
Query(l, r) = prefix(r) − prefix(l−1)
```

The subtraction is `Inverse`. Without it, a BIT cannot answer arbitrary range queries — only prefix queries.

**Why `Scale`:** A range update `[l, r] += v` is stored across two BITs using a difference-array decomposition. The prefix sum formula becomes:

```
prefix(x) = b1.query(x) * x  −  b2.query(x)
```

The `* x` term is position-dependent and cannot be expressed by `Merge` alone, since `Merge` combines two values of the same type — not a value with a plain integer.

Example implementation for integers:

```go
type value int

func (v value) Merge(w internal.MergeableValueType) internal.MergeableValueType {
    return v + w.(value)
}
func (v value) Inverse() internal.MergeableValueType { return -v }
func (v value) Scale(n int) internal.MergeableValueType { return v * value(n) }
```

### `ComparableValue` — required by `NewCmpTree`

```go
type ComparableValue interface {
    internal.ValueType
    Merge(other ComparableValue) ComparableValue  // max(a,b) or min(a,b)
    Identity() ComparableValue                    // −∞ for max, +∞ for min
}
```

Max/min types have no inverse (`max(a,b) − max(a,c)` is not meaningful), so they cannot satisfy `InvertibleValue`. They are confined to the simpler `cmpTree` which uses a single BIT and only answers prefix queries.

Example implementations:

```go
type maxVal int

func (v maxVal) Merge(w ComparableValue) ComparableValue {
    if w.(maxVal) > v { return w }
    return v
}
func (v maxVal) Identity() ComparableValue { return maxVal(math.MinInt64) }

type minVal int

func (v minVal) Merge(w ComparableValue) ComparableValue {
    if w.(minVal) < v { return w }
    return v
}
func (v minVal) Identity() ComparableValue { return minVal(math.MaxInt64) }
```

---

## Operations

### `SumTree` — range update, range query, point update, point query

#### `Modify(i, j KeyType, v MergeableValueType)`

Adds `v` to every element in the closed range `[i, j]`.

Internally decomposes into four point updates across two BITs:

```
b1: +v at l,       −v at r+1
b2: +v*(l−1) at l, −v*r  at r+1
```

This encodes the coefficient and offset of the two-BIT prefix sum formula separately so that `prefix(x)` can reconstruct the correct per-position sum in O(log n).

**Complexity:** O(log n)

#### `Query(i, j KeyType) MergeableValueType`

Returns the sum of all elements in the closed range `[i, j]`.

Computed as:

```
Query(l, r) = prefix(r) − prefix(l−1)

where prefix(x) = b1.query(x) * x − b2.query(x)
```

**Complexity:** O(log n)

#### `Update(key KeyType, value ValueType)`

Point update. Equivalent to `Modify(key, key, value)`. Adds `value` to the single element at `key`.

**Complexity:** O(log n)

#### `Find(key KeyType) ValueType`

Point query. Equivalent to `Query(key, key)`. Returns the accumulated value at the single index `key`.

Derivation: `Query(i, i) = prefix(i) − prefix(i−1)` cancels all elements outside position `i`, leaving only the element's own accumulated value.

**Complexity:** O(log n)

---

### `cmpTree` — monotone point update, prefix query

#### `Update(key KeyType, value ValueType)`

Applies `Merge(current[key], value)` at position `key`. For a max-BIT this means the stored value only ever increases; for a min-BIT it only ever decreases.

**Precondition:** `value` must be "better" than the current value in the Merge ordering — i.e. larger for max, smaller for min. Calling with a worse value is not an error but has no effect, since `Merge` will simply keep the existing value.

**Complexity:** O(log n)

#### `Find(key KeyType) ValueType`

Returns the prefix `Merge` of all elements in `[0..key]` — the maximum (or minimum) value among all indices from 0 up to and including `key`.

Note: this is a **prefix query**, not a point read. The BIT structure distributes each value across multiple ancestor nodes, so reading back a single index is not directly supported. If you need both a point read and a prefix query, maintain a separate flat array alongside the BIT.

**Complexity:** O(log n)

---

## Capability matrix

| Operation | `SumTree` | `cmpTree` (max/min) |
|---|---|---|
| Range update `[l, r]` | ✅ | ❌ no inverse |
| Range query `[l, r]` | ✅ | ❌ no inverse |
| Point update | ✅ | ✅ monotone only |
| Point query (single index) | ✅ | ❌ prefix only |
| Prefix query `[0..i]` | ✅ | ✅ |
| Decrease a max-BIT value | ✅ | ❌ corrupts ancestors |

---

## Two-BIT derivation

A standard BIT answers prefix queries after point updates. To support range updates, the trick is to represent element values implicitly using a difference array, then derive prefix sums of that difference array.

After applying range updates, the value at element `i` is:

```
element[i] = Σ { v : update [l,r] covers i }
```

The prefix sum up to position `x` is:

```
prefix(x) = Σ_{i=1}^{x} element[i]
           = Σ_{i=1}^{x} Σ { v : l ≤ i ≤ r }
```

Rearranging per update `[l, r, v]`, its contribution to `prefix(x)` is:

```
if x < l:        0
if l ≤ x ≤ r:    v * (x − l + 1)  =  v*x  −  v*(l−1)
if x > r:        v * (r − l + 1)  =  v*r  −  v*(l−1)
```

This has the form `coefficient * x − constant`, where the coefficient and constant both change at positions `l` and `r+1`. Two BITs track each part:

```
b1 tracks coefficient changes: +v at l, −v at r+1
b2 tracks constant changes:    +v*(l−1) at l, −v*r at r+1

prefix(x) = b1.query(x) * x  −  b2.query(x)
```

`Scale` computes `v * (l−1)` and `v * r`. `Inverse` computes the negations applied at `r+1` and in the final `Query(l,r) = prefix(r) − prefix(l−1)`.

---

## Complexity summary

| Operation | Time | Space |
|---|---|---|
| Build | O(n) | O(n) — two arrays of size n+2 |
| `Modify` / range update | O(log n) | — |
| `Query` / range query | O(log n) | — |
| `Update` / point update | O(log n) | — |
| `Find` / point or prefix query | O(log n) | — |

Compared to the segment tree with lazy propagation, the Fenwick tree has a smaller constant factor (two tight loops vs recursive descent with pushDown), but is restricted to invertible operations for range queries.
