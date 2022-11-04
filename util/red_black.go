package util

type RedBlackTree struct {
	index               int
	red                 bool
	left, right, parent *RedBlackTree
}

/*
red black tree:
1. Every node is red or black
2. The root must be black
3. Every nil is black
4. Red children must be black
5. The black-height must be the same across the tree.
*/

func (t *RedBlackTree) IsLeftOf(other *RedBlackTree) bool {
	return t.index <= other.index
}

func (t *RedBlackTree) RightRotate() {
	// let's implement this on our own.
	// The goal is to make X Y's right (larger) child, maintaining order.
	x := t      // x is our base
	y := t.left // y is smaller, left

	// If Y has a right child, we must re-parent it to X,
	// which now has a convenient empty left channel
	if y.right != nil {
		x.left = y.right
		x.left.parent = x
	} else {
		x.left = nil
	}

	// Then, we must re-parent X.
	// if X is the root, Y becomes the root.
	if x.parent == nil {
		y.parent = nil
	} else {
		// otherwise, Y becomes the child of whichever side X was on.
		p := x.parent
		y.parent = p

		if p.left == x {
			p.left = y
		} else {
			p.right = y
		}
	}

	// X then becomes Y's right (greater) child
	x.parent = y
	y.right = x
}

func (t *RedBlackTree) LeftRotate() {
	x := t
	y := t.right // Y is larger, right

	// In a left rotate, we put Y as the left parent of X.
	// If left already has a Y, we need to re-parent it.
	// Thankfully, there is a brand-new empty node that perfectly fits it,
	// Larger than X (wouldn't be right of X if smaller),
	// Less than Y (wouldn't be left of Y if larger).
	// Thus, it is correct to place it there.
	if y.left != nil {
		x.right = y.left
		x.right.parent = x
	} else {
		x.right = nil
	}

	// We then need to re-parent X.
	// If it's the root, the swap is trivial.
	if x.parent == nil {
		y.parent = nil
	} else {
		// If not, Y simply replaces where X was.
		p := x.parent
		y.parent = p

		if p.left == x {
			p.left = y
		} else {
			p.right = y
		}
	}

	// set X to left child of Y
	// X is smaller than Y, otherwise, Y wouldn't be the right child.
	x.parent = y
	y.left = x
}

func (t *RedBlackTree) Insert(index int) *RedBlackTree {
	newLeaf := &RedBlackTree{
		index: index,
		red:   true, // default to red because it's the safest algorithm.
	}

	if t == nil {
		newLeaf.red = false
		return newLeaf
	}

	// First, we must locate the initial location.
	p := t
	for {
		if newLeaf.IsLeftOf(p) { // less/same = left
			if p.left == nil {
				break
			}
			p = p.left
		} else { // greater = right
			if p.right == nil {
				break
			}
			p = p.right
		}
	}

	newLeaf.parent = p
	if index <= p.index {
		p.left = newLeaf
	} else {
		p.right = newLeaf
	}

	if p.red { // red cannot parent red. Fix the insert, as it's going to increase the black height.
		return t.Repair(newLeaf)
	}

	return t
}

func (t *RedBlackTree) isRed() bool {
	return t != nil && t.red
}

func (t *RedBlackTree) Repair(newLeaf *RedBlackTree) *RedBlackTree {
	var p, gp, uncle *RedBlackTree
	//redo:
	p = newLeaf.parent
	if p == nil { // root must be black
		newLeaf.red = false
		return newLeaf
	}
	gp = p.parent
	if gp != nil {
		uncle = Ternary(newLeaf.IsLeftOf(gp), gp.right, gp.left)
	}

	fixRoot := func() {
		for t.parent != nil {
			t = t.parent
		}
	}

	if newLeaf.isRed() && p.isRed() {
		if uncle.isRed() {
			// If the uncle is red, we can maintain the black level by pushing it down the tree.
			p.red = false
			uncle.red = false
			gp.red = true

			fixRoot()
			return t.Repair(gp)
		} else {
			if newLeaf.IsLeftOf(p) == p.IsLeftOf(gp) {
				// Rotating the grandparent can resolve the double-red.
				if newLeaf.IsLeftOf(gp) {
					gp.RightRotate()
				} else {
					gp.LeftRotate()
				}

				// Make sure we re-color it.
				gp.red = true
				p.red = false
				newLeaf.red = true
				fixRoot()
				return t.Repair(p)
			} else {
				// Something of a diamond. Make the new child the new parent.
				if newLeaf.IsLeftOf(p) {
					p.RightRotate()
				} else {
					p.LeftRotate()
				}

				fixRoot()
				return t.Repair(p) // Treat the old parent like the new child, because it's now two deep, and red.
			}
		}
	}

	return t
}
