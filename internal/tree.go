package internal

type Tree interface {
	// Find finds the node and return its value
	Find(key KeyType) ValueType

	// Update the value associated with a key inside the tree
	Update(key KeyType, value ValueType)
}

type BST interface {
	Tree

	// Empty checks whether the tree is empty
	Empty() bool

	// Size returns the size of the tree
	Size() int

	// Clear destroys the tree
	Clear()

	// Finds the node and return it as a iterator
	FindIt(key KeyType) Iterator

	// Creates the node iterator start from most left node
	Iterator() Iterator

	// Insert inserts the key-value pair into the tree
	Insert(key KeyType, value ValueType)

	// Print the tree by traversal order
	Preorder()

	// Delete deletes the node by key
	Delete(key KeyType)

	// Min finds min element in the tree
	Min() ValueType

	// Max finds min element in the tree
	Max() ValueType
}

type RangeQuery interface {
	Tree

	// Query returns the merged value over the closed range [i, j].
	Query(i, j KeyType) MergeableValue

	// Modify performs a update applicable over the closed range [i, j] on the value v.
	Modify(i, j KeyType, v MergeableValue)
}
