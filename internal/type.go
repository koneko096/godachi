package internal

// KeyType must implement comparison operations
type KeyType interface {
	// Key must be compared against order
	LessThan(KeyType) bool

	// Key must be exactly matched
	Equal(KeyType) bool
}

type ValueType interface{}

// Iterator is tree node pointer which can traverse the tree
type Iterator interface {
	// Check if the iterator is nil node
	IsNil() bool

	// Advance the iterator to next element
	Next() Iterator

	// Obtain the key of node
	Key() KeyType

	// Obtain the value of node
	Value() ValueType

	// Update the value of the node
	Set(v ValueType)
}
