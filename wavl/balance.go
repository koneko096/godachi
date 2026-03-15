package wavl

// ---- WAVL rebalance: Insert ----------------------------------------------
//
// After inserting a leaf, walk up from the parent fixing 0-children.
//
//	Case 1 — 0,1-node: promote cur, continue upward.
//	Case 2 — 0,2-node: rotate (single or double), always terminates.
//	  single: demote cur (old root goes down one level)
//	  double: demote cur, demote the other child, promote newRoot
func rebalanceInsertUp(cur *node) *node {
	var last *node
	for cur != nil {
		last = cur
		ldiff := rankDiff(cur, cur.left)
		rdiff := rankDiff(cur, cur.right)

		if ldiff >= 1 && rdiff >= 1 {
			break // valid WAVL node, done
		}

		if ldiff+rdiff == 1 {
			// Case 1: 0,1-node → promote and continue upward
			cur.promote()
			cur = cur.parent
			continue
		}

		// Case 2: 0,2-node → rotate, always terminates
		// WAVL rank rules (Haeupler, Sen, Tarjan):
		//   single: demote p (cur)
		//   double: demote p (cur), demote x (child), promote y (newRoot)
		// capture x BEFORE rotation — x is the deficient child of cur
		var x *node
		if rankDiff(cur, cur.left) < rankDiff(cur, cur.right) {
			x = cur.left // left is 0-child
		} else {
			x = cur.right // right is 0-child
		}
		newRoot, double := cur.fixRotation()
		cur.demote() // always demote p
		if double {
			x.demote()        // demote x (went down one level)
			newRoot.promote() // promote y (came up two levels)
		}
		cur = newRoot
		break
	}
	if cur != nil {
		return cur.treeRoot()
	}
	return last.treeRoot()
}

// ---- WAVL rebalance: Delete ----------------------------------------------
//
// After deletion, walk up from the deletion point fixing 3-children
// and 2,2-leaves.
//
//	2,2-leaf:          demote leaf, continue upward.
//	Case 1 — 3-child, sibling is 2,2: demote cur + sibling, continue upward.
//	Case 2 — 3-child, sibling has a 1-child: rotate, always terminates.
//	  single: demote cur twice, promote newRoot
//	  double: demote cur twice, demote old child, promote newRoot
//	Case 3 — 3-child, sibling diff == 2: demote cur, continue upward.
func rebalanceDeleteUp(cur *node) *node {
	var last *node
	for cur != nil {
		last = cur
		ldiff := rankDiff(cur, cur.left)
		rdiff := rankDiff(cur, cur.right)

		// 2,2-leaf: rank too high after deletion, demote and continue
		if cur.left == nil && cur.right == nil && cur.rank > 1 {
			cur.demote()
			cur = cur.parent
			continue
		}

		if ldiff <= 2 && rdiff <= 2 {
			break // valid WAVL node, done
		}

		// 3-child exists — sibling is on the opposite side
		sibling := cur.right
		if rdiff > 2 {
			sibling = cur.left
		}

		// nil sibling guard — should not happen in valid WAVL tree
		if sibling == nil {
			cur.demote()
			cur = cur.parent
			continue
		}

		slDiff := rankDiff(sibling, sibling.left)
		srDiff := rankDiff(sibling, sibling.right)

		if slDiff == 2 && srDiff == 2 {
			// Case 1: sibling is 2,2-node → demote both, continue upward
			cur.demote()
			sibling.demote()
			cur = cur.parent
			continue
		}

		// Case 2: sibling has at least one 1-child → rotate, always terminates
		// WAVL rank rules (Haeupler, Sen, Tarjan):
		//   single: demote p (cur) once,  promote s (newRoot)
		//   double: demote p (cur) twice, demote s, promote y (newRoot)
		// sibling s is already captured above — use it directly
		newRoot, double := cur.fixRotation()
		cur.demote() // always demote p at least once
		if double {
			cur.demote()     // second demote for p
			sibling.demote() // demote s
		}
		newRoot.promote() // always promote newRoot (s for single, y for double)
		cur = newRoot
		break
	}
	if cur != nil {
		return cur.treeRoot()
	}
	return last.treeRoot()
}

// fixRotation performs a single or double rotation to resolve structural
// imbalance at n. Uses rank diff to decide rotation direction (which side
// is deficient), and structural height to decide single vs double rotation.
// All rank management is handled by the caller.
// Returns (newRoot, wasDouble).
func (n *node) fixRotation() (*node, bool) {
	// capture grandparent BEFORE any rotation — n is still gp's direct child
	gp := n.parent
	isLeftChild := gp != nil && gp.left == n

	// rank diff decides DIRECTION: rotate toward the deficient side
	// (0-child on insert, 3-child on delete)
	ldiff := rankDiff(n, n.left)
	rdiff := rankDiff(n, n.right)

	var newRoot *node
	double := false

	if ldiff < rdiff {
		// left side is deficient — bring left child up (rotate right)
		nl := n.left
		if nl.right.height() > nl.left.height() {
			// left-right: double rotation
			double = true
			n.left = nl.right.rotateLeft(nl)
			nl = n.left
		}
		newRoot = nl.rotateRight(n)
	} else {
		// right side is deficient — bring right child up (rotate left)
		nr := n.right
		if nr.left.height() > nr.right.height() {
			// right-left: double rotation
			double = true
			n.right = nr.left.rotateRight(nr)
			nr = n.right
		}
		newRoot = nr.rotateLeft(n)
	}

	// fix grandparent link
	newRoot.parent = gp
	if gp != nil {
		if isLeftChild {
			gp.left = newRoot
		} else {
			gp.right = newRoot
		}
	}

	return newRoot, double
}

func (n *node) rotateRight(p *node) *node {
	p.insertLeft(n.right)
	n.parent = p.parent
	n.insertRight(p)
	return n
}

func (n *node) rotateLeft(p *node) *node {
	p.insertRight(n.left)
	n.parent = p.parent
	n.insertLeft(p)
	return n
}
