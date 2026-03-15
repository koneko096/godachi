# Weak AVL Tree (WAVL Tree)
A WAVL tree combines the AVL tree's more balanced nature with Red-Black tree efficiency by using **ranks** instead of heights to maintain balance. The key relaxation: rank difference between any parent and child must be 1 or 2 — allowing the 2,2-node that AVL trees forbid, while remaining stricter than Red-Black trees.

The result is a tree that:
- Without deletions, behaves **exactly** like an AVL tree
- With deletions, height stays at most that of an AVL tree with the same insertions but no deletions
- Always performs **at most 2 rotations** per insert or delete (AVL delete requires O(log n))

---

## Core Concept

### Rank
Each node carries an integer **rank** — an approximation of its distance to the farthest leaf. Unlike AVL trees where rank equals height exactly, WAVL ranks can diverge from height after deletions.

**Three invariants must hold at all times:**

| Rule | Requirement |
|---|---|
| External-Node Property | Every external (null) node has rank **-1** |
| Rank-Difference Property | Every non-root node's rank difference with its parent is **1 or 2** |
| Leaf Property | A leaf (node with two null children) must have rank exactly **0** |

The **rank difference** of a node x is `rank(parent(x)) - rank(x)`. A node is called an **i,j-node** where i is its left child's rank difference and j is its right child's rank difference.

**Valid node types:**

```
1,1-node  →  both children are 1-children   (e.g. internal node with two rank-0 leaves)
1,2-node  →  rank diffs are 1 and 2         (normal internal node)
2,2-node  →  both children are 2-children   ← only in WAVL, not AVL
```

The 2,2-node is the key distinction from AVL trees. AVL trees forbid it because it would mean two children with equal height under a parent with rank = height + 2. WAVL allows it, which avoids expensive rebalancing cascades on deletion.

**Rank convention used in this implementation:**

```
null node  → rank -1
leaf node  → rank  0   (inserted with rank 0, diffs to null children = 1 each → valid 1,1-node)
internal   → rank ≥ 1  (promoted as needed during rebalancing)
```

This is the Haeupler/Sen/Tarjan (2015) formulation. Some texts use leaf=1 with null=-1 — the
invariants are equivalent, just shifted by one.

**Height bound:** A WAVL tree with n nodes has height at most **2 log₂(n+1)**, same asymptotic bound as AVL and Red-Black trees.

**Encoding ranks efficiently:** Since only rank differences of 1 or 2 are valid, ranks can be stored as bit flags rather than integers:

```
Option 1: Two bits per node
  bit 0 = rank diff to left child  (0 → diff is 1,  1 → diff is 2)
  bit 1 = rank diff to right child (0 → diff is 1,  1 → diff is 2)

Option 2: Single parity bit per node
  store parity (even/odd) of absolute rank
  same parity as parent → rank diff is 2
  different parity      → rank diff is 1
```

### Rotations
WAVL uses the same left/right rotations as AVL and Red-Black trees, but **rank adjustments accompany every rotation**. At most **2 rotations** are ever needed per operation.

Importantly, rotation direction is determined by **rank diff** (which side has the violation), and single vs double is determined by **structural height** of the grandchildren. These are two separate decisions.

**Single rotation (e.g. left rotation at p):**
```
    p                s
   / \              / \
  A   s     →      p   C
     / \          / \
    B   C        A   B

Rank updates:
  rank(p) -= 1  (demote — p goes down one level)
  rank(s) += 1  (promote — s comes up one level)
```

**Double rotation (right-left at p):**
```
    p                  y
   / \               /   \
  A   s     →       p     s
     / \           / \   / \
    y   C         A   B1 B2  C
   / \
  B1  B2

Rank updates:
  rank(p) -= 1  (demote)
  rank(s) -= 1  (demote)
  rank(y) += 1  (promote — y came up two levels)
```

The rank adjustments are what distinguish WAVL rotations from naive BST rotations — the structural move alone is not enough.

---

## Operations

### Find
Identical to any BST search — no rank logic needed for reads.

```
find(t, key):
  if t is null: return null
  if key == t.key: return t
  if key  < t.key: return find(t.left, key)
  else:            return find(t.right, key)
```

Time complexity: **O(log n)** worst case (height bound of 2 log n).

### Insert
Insert as a standard BST leaf with rank 0, then walk up fixing rank violations.

```
insert(root, key):
  1. BST insert — place new node q as a leaf with rank 0
  2. Walk up from q's parent p, fixing violations bottom-up
```

A fresh leaf at rank 0 has null children at rank -1, giving rank diffs of 1 each — a valid 1,1-node. No immediate violation. A violation only appears at p when p's rank equals q's rank (diff = 0), which happens if p was previously promoted to match q's subtree height.

**Rebalancing cases** (at node p, child q is a 0-child — rank diff = 0):

**Case 1: p is a 0,1-node (q's sibling has rank-difference 1)**
```
Action: Promote p  (rank(p) += 1)
Effect: Fixes the 0-child violation at p.
        May create a new violation at p's parent → continue upward.

Before:          After:
  p[r]            p[r+1]
  / \             / \
q[r] s[r-1]    q[r] s[r-1]
(diff=0)(diff=1)  (diff=1)(diff=2) ✓
```

**Case 2: p is a 0,2-node (q's sibling has rank-difference 2)**

Sub-case 2a — single rotation (q is a 1,2 or 2,1-node, outer grandchild is heavier):
```
Action: Single rotation + demote p
Effect: Fixes violation in O(1). Stop — no further propagation.

Before:            After (rotate right at p):
    p[r]                q[r]
   / \                 / \
  q[r] s[r-2]  →    A    p[r-1]
  / \                    / \
 A   B                  B   s[r-2]
(outer A is heavier)
```

Sub-case 2b — double rotation (q's inner grandchild y is heavier):
```
Action: Double rotation + demote p + demote q + promote y
Effect: Fixes violation in O(1). Stop — no further propagation.

Before:              After (double rotation):
    p[r]                  y[r]
   / \                  /     \
  q[r] s[r-2]  →      q[r-1]  p[r-1]
  / \                  / \    / \
 A   y[r-1]           A   B1 B2  s[r-2]
    / \
   B1  B2
```

Key property: **insertion never creates a 2,2-node**, so a tree built only by insertions is always a valid AVL tree.

Time complexity: **O(log n)** for the BST traversal, **O(1) amortized** for rebalancing (at most 2 rotations, O(log n) promotes).

### Delete
Delete using standard BST removal, then walk up fixing rank violations. This is where WAVL's advantage over AVL is most visible — AVL deletion can cascade O(log n) rotations, WAVL always does **at most 2**.

```
delete(root, key):
  1. BST delete — if node has two children, copy successor's key/value
     into the node, then delete the successor (which has at most one child)
  2. Walk up from the deleted node's parent, fixing violations bottom-up
```

After deletion, two kinds of violation can appear:

- **3-child**: a node whose rank diff with its parent is 3 (parent rank too high)
- **2,2-leaf**: a leaf whose rank is > 0 (rank was not decremented when child was removed)

**Rebalancing cases** (at node p, one side is a 3-child):

**2,2-leaf special case:**
```
Action: Demote the leaf (rank -= 1)
Effect: May propagate upward → continue.
```

**Case 1: sibling s has rank-difference 1, and s is a 2,2-node**
```
Action: Demote p, demote s
Effect: May propagate violation upward → continue.
```

**Case 2: sibling s has rank-difference 1, and s is NOT a 2,2-node**

Sub-case 2a — single rotation (s's far child is a 1-child):
```
Action: Single rotation at p + demote p + promote s (newRoot)
Effect: Terminates. O(1) rotations.
```

Sub-case 2b — double rotation (s's near child y is a 1-child, far child is a 2-child):
```
Action: Double rotation + demote p + demote s + promote y (newRoot)
Effect: Terminates. O(1) rotations.
```

**Case 3: sibling s has rank-difference 2**
```
Action: Demote p only
Effect: May propagate upward → continue.
```

The crucial insight: once a rotation fires (Cases 2a or 2b), rebalancing **always terminates** — no further propagation. Only demotions (Cases 1 and 3, and the 2,2-leaf case) can propagate.

Time complexity: **O(log n)** for the BST traversal, **O(1) amortized** for rebalancing.

---

## Implementation Notes

### Rotation Direction vs Single/Double Decision

Two separate decisions are made during `fixRotation`:

```
1. Direction  → which side to rotate toward
               determined by RANK DIFF (which side has the violation)
               ldiff < rdiff → left side is deficient → rotate right (bring left child up)
               ldiff > rdiff → right side is deficient → rotate left (bring right child up)

2. Single vs double → whether a pre-rotation is needed
               determined by HEIGHT of grandchildren (structural, not rank)
               if inner grandchild is taller → double rotation needed
               otherwise → single rotation sufficient
```

Using height (not rank) for the single/double decision is essential — stored ranks can diverge from height after deletions and promotions, so using rank would pick the wrong rotation type.

### Rank Management After Rotation

Rotations themselves make no rank changes. All rank adjustments are applied by the caller after `fixRotation` returns:

```
Insert single:  demote p
Insert double:  demote p, demote x (0-child), promote y (newRoot)

Delete single:  demote p, promote s (newRoot)
Delete double:  demote p, demote s, promote y (newRoot)
```

---

## Comparison with AVL and Red-Black Trees

| Property | AVL Tree | **WAVL Tree** | Red-Black Tree |
|---|---|---|---|
| Balance rule | Height diff ≤ 1 | Rank diff = 1 or 2 | Color constraints |
| Height bound | 1.44 log n | 2 log n | 2 log n |
| Insert rotations | ≤ 2 | ≤ 2 | ≤ 2 |
| Delete rotations | **O(log n)** | **≤ 2** | ≤ 3 |
| Insert-only behavior | AVL | **Exactly AVL** | Red-Black |
| 2,2-nodes allowed | ❌ | ✅ | N/A |
| Initial leaf rank | height = 1 | **0** (null children = -1, diffs = 1,1) | N/A (color-based) |
| Implementation complexity | High | Moderate | Moderate |

WAVL strictly sits between AVL and Red-Black: every AVL tree is a valid WAVL tree (with ranks equal to heights), and every WAVL tree can be recolored into a valid Red-Black tree — but not all Red-Black trees come from WAVL trees.
