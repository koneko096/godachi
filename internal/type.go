package internal

// KeyType must implement comparison operations
type KeyType interface {
	// Key must be compared against order
	LessThan(KeyType) bool

	// Key must be exactly matched
	Equal(KeyType) bool
}

type ValueType interface{}

// MergeableValue is implemented by any value that participates in
// range-aggregate queries and lazy range modifications.
//
// The three methods encode the full monoid semantics for a given operation:
//
//   - Merge   – associatively combines two aggregated values (left-to-right).
//     Used by pullUp to build a subtree aggregate from its children.
//
//   - Apply   – returns the new aggregate after this tag is applied uniformly
//     to a subtree whose current aggregate is agg and whose node
//     count is sz.
//     Examples:
//     range-add / sum:    agg + tag*sz
//     range-assign / min: tag   (size irrelevant)
//
//   - Compose – returns the single tag equivalent to first applying the
//     receiver and then applying newer.  Used to stack lazy tags
//     without eager propagation.
//     Examples:
//     additive tags:   self + newer
//     assign tags:     newer  (later assignment wins)
type MergeableValue interface {
	Merge(other MergeableValue) MergeableValue
	Apply(agg MergeableValue, sz int) MergeableValue
	Compose(newer MergeableValue) MergeableValue
}

// InvertibleValue extends MergeableValue for sum-based range update +
// range query. Requires additive inverse and scalar multiply.
// Not usable for max/min (no inverse exists).
type InvertibleValue interface {
	MergeableValue

	Inverse() MergeableValue
	Scale(n int) MergeableValue
}

// ComparableValue extends MergeableValue for monotone point update + prefix query.
// Merge must be an idempotent, associative operation (max or min).
// Update is only valid when the new value is strictly "better" than the old
// one in the Merge ordering (i.e. monotone — never decrease for max-BIT,
// never increase for min-BIT).
type ComparableValue interface {
	MergeableValue

	// Identity returns the identity element for Merge.
	// For max-BIT: −∞. For min-BIT: +∞.
	Identity() ComparableValue
}

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
