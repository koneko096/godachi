# Treap

Treap comes from **tree + heap** — it is a BST (Binary Search Tree) that simultaneously satisfies the heap property using randomly assigned priorities. This dual invariant gives O(log n) expected height without any explicit rebalancing.

## Core Concept

### Node Structure

Each node holds two keys:
- **key** — determines BST ordering (in-order traversal gives sorted keys)
- **priority** — assigned randomly at creation, determines heap ordering (parent's priority ≥ children's)

Because priorities are random and independent of keys, the resulting tree has the same distribution as a random BST, giving **O(log n) expected height**.

```
Node structure:
┌──────────────┐
│  key = 5     │  ← BST property: left.key < 5 < right.key
│  priority=91 │  ← Heap property: 91 ≥ children's priorities
│  value       │
│  left, right │
└──────────────┘
```

No adversarial input can degenerate the tree — the priorities are random, not derived from the data.

### Split

`splitLeft(t, k)` and `splitRight(t, k)` divide the treap into two treaps by **key value**:

```
splitLeft(t, k)  → L: keys ≤ k  |  R: keys > k
splitRight(t, k) → L: keys < k  |  R: keys ≥ k
```

```
splitLeft(t, k):
  if t is null: return (null, null)
  propagate(t)                         // flush lazy tags first
  if k < t.key:
    (ll, lr) = splitLeft(t.left, k)
    t.left   = lr
    update(t)
    return (ll, t)                     // t goes to R
  else:
    (rl, rr) = splitLeft(t.right, k)
    t.right  = rl
    update(t)
    return (t, rr)                     // t goes to L
```

**Example** — `splitLeft([1,2,3,4,5], k=3)`:
```
         3(p=9)                 L=[1,2,3]   R=[4,5]
        /      \      →
    2(p=7)   4(p=3)
    /            \
 1(p=2)        5(p=1)
```

Time complexity: **O(log n)** expected.

### Merge

`merge(L, R)` combines two treaps where every key in L is less than every key in R, producing a single valid treap.

The node with the higher priority becomes the new root — this preserves the heap property by construction.

```
merge(L, R):
  if L is null: return R
  if R is null: return L
  propagate(L)
  propagate(R)
  if L.priority > R.priority:
    L.right = merge(L.right, R)
    update(L)
    return L
  else:
    R.left = merge(L, R.left)
    update(R)
    return R
```

Time complexity: **O(log n)** expected.

---

## Operations

### Find

Standard BST search — follow left/right based on key comparison.

```
find(t, k):
  if t is null:    return null
  if k == t.key:   return t
  if k  < t.key:   return find(t.left,  k)
  else:             return find(t.right, k)
```

Time complexity: **O(log n)** expected.

### Insert

Split at the target key, create a new single-node treap, then merge everything back.

```
insert(root, k, val):
  new_node     = Node(k, val, random_priority)
  (L, R)       = splitRight(root, k)   // L: keys < k  |  R: keys ≥ k
  root         = merge(merge(L, new_node), R)
```

If `k` already exists it is removed first so there are no duplicate keys.

Time complexity: **O(log n)** expected.

### Delete

Split twice to isolate all nodes with the target key, discard them, then merge the remaining two parts.

```
delete(root, k):
  (L, mid_R) = splitRight(root, k)   // L: keys < k
  (mid, R)   = splitLeft(mid_R, k)   // mid: keys == k  |  R: keys > k
  // discard mid
  root = merge(L, R)
```

Time complexity: **O(log n)** expected.

---

## Range Query (lazy propagation)

`rangeTree` extends the treap with two additional operations over closed key intervals `[lo, hi]`:

```
Node structure (range tree):
┌──────────────┐
│  key = 5     │
│  priority=91 │
│  value       │  ← point value (MergeableValue)
│  agg         │  ← cached aggregate of entire subtree
│  lazy        │  ← pending tag to push to children
│  left, right │
└──────────────┘
```

### MergeableValue

All monoid semantics are encoded by the value type — the tree carries no operation-specific configuration.

```
interface MergeableValue:
  Merge(other)       → combine two aggregates (must be associative)
  Apply(agg, size)   → apply this tag to a subtree aggregate of given size
  Compose(newer)     → combine two pending tags into one
```

Two example configurations:

| | range-add / sum | range-assign / min |
|---|---|---|
| `Merge(a,b)` | `a + b` | `min(a, b)` |
| `Apply(agg, sz)` | `agg + tag×sz` | `tag` |
| `Compose(newer)` | `self + newer` | `newer` |

### Lazy Propagation

Every structural operation follows the same two-step discipline:

```
propagate(t)   // before descending:  flush lazy tag to children
...            // structural work
update(t)      // after ascending:    recompute sz and agg
```

`propagate` stamps the pending tag onto each child in O(1) via `applyTag`, then clears it:

```
propagate(t):
  if t.lazy is null: return
  applyTag(t.left,  t.lazy)
  applyTag(t.right, t.lazy)
  t.lazy = null

applyTag(t, tag):
  t.value = tag.Apply(t.value, 1)
  t.agg   = tag.Apply(t.agg, t.sz)
  t.lazy  = t.lazy == null ? tag : t.lazy.Compose(tag)
```

`update` rebuilds `agg` from the node's own value and its children's aggregates:

```
update(t):
  t.sz  = 1 + size(t.left) + size(t.right)
  t.agg = t.value
  if t.left  != null: t.agg = t.left.agg.Merge(t.agg)
  if t.right != null: t.agg = t.agg.Merge(t.right.agg)
```

### Query

Returns the aggregate over all values whose keys fall in `[lo, hi]`.

```
query(root, lo, hi):
  (l, mr) = splitRight(root, lo)   // l:  keys < lo
  (m, r)  = splitLeft(mr, hi)      // m:  keys in [lo, hi]  |  r: keys > hi
  result  = m.agg                  // aggregate already cached at subtree root
  root    = merge(l, merge(m, r))
  return result
```

Time complexity: **O(log n)** expected.

### Modify

Applies a tag uniformly to every node whose key falls in `[lo, hi]`.

```
modify(root, lo, hi, tag):
  (l, mr) = splitRight(root, lo)
  (m, r)  = splitLeft(mr, hi)
  applyTag(m, tag)                 // O(1) — stamps tag on subtree root only
  root    = merge(l, merge(m, r))  // propagate flushes tag lazily on descent
```

Time complexity: **O(log n)** expected.
