package treap

import (
	"testing"

	"github.com/koneko096/godachi/internal"
)

func TestPreorder(t *testing.T) {
	var tree internal.BST = NewTree()

	tree.Insert(key(1), "123")
	tree.Insert(key(3), "234")
	tree.Insert(key(4), "dfa3")
	tree.Insert(key(6), "sd4")
	tree.Insert(key(5), "jcd4")
	tree.Insert(key(2), "bcd4")
	if tree.Size() != 6 {
		t.Error("Error size")
	}
	tree.Preorder()
}

func TestFindElement(t *testing.T) {
	var tree internal.BST = NewTree()

	tree.Insert(key(1), "123")
	tree.Insert(key(3), "234")
	tree.Insert(key(4), "dfa3")
	tree.Insert(key(6), "sd4")
	tree.Insert(key(5), "jcd4")
	tree.Insert(key(2), "bcd4")
	tree.Preorder()

	n := tree.FindIt(key(4))
	if n == nil {
		t.Error("Error value not found")
		return
	}
	if n.Value() != "dfa3" {
		t.Error("Error value")
	}
	n.Set("bdsf")
	if n.Value() != "bdsf" {
		t.Error("Error value modify")
	}
	value := tree.Find(key(5)).(string)
	if value != "jcd4" {
		t.Error("Error value after modifyed other node")
	}
}

func TestIterator(t *testing.T) {
	var tree internal.BST = NewTree()

	tree.Insert(key(1), "123")
	tree.Insert(key(3), "234")
	tree.Insert(key(4), "dfa3")
	tree.Insert(key(6), "sd4")
	tree.Insert(key(5), "jcd4")
	tree.Insert(key(2), "bcd4")

	it := tree.Iterator()

	id := 1
	for it.IsNil() {
		if it.Key() != key(id) {
			t.Error("Iterator not ordered")
		}

		id++
		it = it.Next()
	}

}

func TestDelete(t *testing.T) {
	var tree internal.BST = NewTree()

	tree.Insert(key(1), "123")
	tree.Insert(key(3), "234")
	tree.Insert(key(4), "dfa3")
	tree.Insert(key(6), "sd4")
	tree.Insert(key(5), "jcd4")
	tree.Insert(key(2), "bcd4")

	for i := 1; i <= 6; i++ {
		tree.Delete(key(i))
		if tree.Size() != 6-i {
			t.Error("Delete Error")
		}
	}

	for i := 1; i <= 6; i++ {
		if tree.Find(key(i)) != nil {
			t.Error("Element not deleted")
			break
		}
	}
}

func TestClear(t *testing.T) {
	var tree internal.BST = NewTree()

	tree.Insert(key(1), "bcd4")
	tree.Clear()
	tree.Preorder()
	if tree.Find(key(1)) != nil {
		t.Error("Can't clear")
	}
}

func TestDelete3(t *testing.T) {
	var tree internal.BST = NewTree()
	tree.Insert(key(4), "1qa")
	tree.Insert(key(2), "2ws")
	tree.Insert(key(3), "3ed")
	tree.Insert(key(1), "4rf")
	tree.Insert(key(8), "5tg")
	tree.Insert(key(5), "6yh")
	tree.Insert(key(7), "7uj")
	tree.Insert(key(9), "8ik")
	tree.Delete(key(1))
	tree.Delete(key(2))

	if tree.Find(key(2)) != nil {
		t.Error("Element not deleted")
	}
	if tree.Find(key(5)) == nil {
		t.Error("Element not existed")
	}
}

func TestDelete2(t *testing.T) {
	var tree internal.BST = NewTree()
	tree.Insert(key(5), "1qa")
	tree.Insert(key(3), "2ws")
	tree.Insert(key(8), "3ed")
	tree.Insert(key(2), "4rf")
	tree.Insert(key(4), "5tg")
	tree.Insert(key(7), "6yh")
	tree.Insert(key(9), "7uj")
	tree.Insert(key(1), "8ik")
	tree.Delete(key(9))
	tree.Delete(key(6))

	if tree.Find(key(6)) != nil {
		t.Error("Element not deleted")
	}
	if tree.Find(key(5)) == nil {
		t.Error("Element not existed")
	}
}
