package treap

import (
	"testing"

	"github.com/koneko096/godachi/internal"
)

func TestAdd(t *testing.T) {
	var tree internal.RangeQuery = NewTree()

	tree.Insert(key(1), "123")
	tree.Insert(key(3), "234")
	tree.Insert(key(4), "dfa3")
	tree.Insert(key(6), "sd4")
	tree.Insert(key(5), "jcd4")
	tree.Insert(key(2), "bcd4")
}
