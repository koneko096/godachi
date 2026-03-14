package internal

type Tree interface {
	// Find finds the node and return its value
	Find(key KeyType) ValueType

	// Empty checks whether the tree is empty
	Empty() bool

	// Size returns the size of the tree
	Size() int

	// Clear destroys the tree
	Clear()

	// Insert inserts the key-value pair into the tree
	Insert(key KeyType, value ValueType)

	// Delete deletes the node by key
	Delete(key KeyType)

	// Min finds min element in the tree
	Min() ValueType

	// Max finds min element in the tree
	Max() ValueType
}

type BST interface {
	Tree

	// Finds the node and return it as a iterator
	FindIt(key KeyType) Iterator

	// Creates the node iterator start from most left node
	Iterator() Iterator

	// Print the tree by traversal order
	Preorder()
}

type RangeQuery interface {
	Tree

	// Alter(i, j int, v ValueType)

	// Query(i, j int) ValueType
}
