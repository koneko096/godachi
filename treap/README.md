# Treap
Treap comes from **tree + heap** — it is a BST (Binary Search Tree) that simultaneously satisfies the heap property using randomly assigned priorities. This dual invariant gives O(log n) expected height without any explicit rebalancing.

## Core Concept

### Priority
Each node holds two values:
- **key** — determines BST ordering (in-order traversal gives sorted keys)
- **priority** — assigned randomly at creation, determines heap ordering (parent's priority ≥ children's)

Because priorities are random and independent of keys, the resulting tree has the same distribution as a random BST, giving **O(log n) expected height**.

```
Node structure:
┌──────────────┐
│  key = 5     │  ← BST property: left < 5 < right
│  priority=91 │  ← Heap property: 91 ≥ children's priorities
│  left, right │
└──────────────┘
```

No adversarial input can degenerate the tree — the priorities are random, not derived from the data.

### Split
`split(t, k)` divides the treap into two treaps:
- **L**: the first k nodes in in-order traversal
- **R**: the remaining nodes

The key is maintained implicitly via subtree **size** fields — no key is stored.

```
split(t, k):
  if t is null: return (null, null)
  push_down(t)                         // flush lazy tags first
  left_size = size(t->left)

  if left_size >= k:
    (ll, lr) = split(t->left, k)
    t->left  = lr
    update(t)
    return (ll, t)
  else:
    (rl, rr) = split(t->right, k - left_size - 1)
    t->right = rl
    update(t)
    return (t, rr)
```

**Example** — split `[A,B,C,D,E]` at k=2:
```
        D(p=9)                     L=[A,B]   R=[C,D,E]
       /      \        →
    B(p=7)   E(p=3)
   /    \
 A(p=2) C(p=5)
```

Time complexity: **O(log n)** expected.

### Merge
`merge(L, R)` combines two treaps where every node in L precedes every node in R (in array order), producing a single valid treap.

The node with the higher priority becomes the new root — this preserves the heap property by construction.

```
merge(L, R):
  if L is null: return R
  if R is null: return L

  if L->priority > R->priority:
    L->right = merge(L->right, R)
    update(L)
    return L
  else:
    R->left = merge(L, R->left)
    update(R)
    return R
```

**Example** — merge `[A,B,C]` and `[D,E]` where D has highest priority:
```
[A,B,C]  +  [D,E]   →       D(p=9)
                            /      \
                         [A,B,C]   E(p=3)
```

No rotations needed — the heap property drives the correct structure automatically.

Time complexity: **O(log n)** expected.

---

## Operations

### Find
Standard BST search — follow left/right based on key comparison.

```
find(t, key):
  if t is null: return null
  if key == t->key: return t
  if key  < t->key: return find(t->left, key)
  else:             return find(t->right, key)
```

Time complexity: **O(log n)** expected.

### Insert
Split at the target position, create a new single-node treap, then merge everything back.

```
insert(root, i, val):
  new_node = Node(val, random_priority, size=1)
  (L, R)   = split(root, i)
  root     = merge(merge(L, new_node), R)
```

**Example** — insert X at index 2 in `[A,B,C,D,E]`:
```
split(root, 2)        →  L=[A,B],     R=[C,D,E]
merge(L, X)           →  [A,B,X]
merge([A,B,X], R)     →  [A,B,X,C,D,E]  ✓
```

No rebalancing step — merge naturally places the new node where its random priority dictates.

Time complexity: **O(log n)** expected.

### Delete
Split twice to isolate the target node, discard it, then merge the remaining two parts.

```
delete(root, i):
  (L, mid_R) = split(root, i)      // [0..i-1]  |  [i..n-1]
  (mid, R)   = split(mid_R, 1)     // [i]        |  [i+1..n-1]
  // mid is now a single node — discard it
  root = merge(L, R)
```

**Example** — delete index 2 (C) from `[A,B,C,D,E]`:
```
split(root, 2)    →  L=[A,B],   mid_R=[C,D,E]
split(mid_R, 1)   →  mid=[C],   R=[D,E]
discard mid
merge(L, R)       →  [A,B,D,E]  ✓
```

Time complexity: **O(log n)** expected.
