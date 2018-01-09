package HeraDB

import "testing"
import "fmt"

func TestCreateBTree(t *testing.T) {
	n := CreateBTree()
	n.Put([]byte("baz"))
	n.Put([]byte("foo"))
	n.Put([]byte("bar"))
	n.Put([]byte("foo"))
	if len(n.inodes) != 3 {
		t.Fatalf("exp=3; got=%d", len(n.inodes))
		fmt.Printf("%s", n)
	}
	if k := n.inodes[0].key; string(k) != "bar" {
		t.Fatalf("exp=<bar>; got=<%s>", k)
	}
	if k := n.inodes[1].key; string(k) != "baz" {
		t.Fatalf("exp=<baz>; got=<%s>", k)
	}
	if k := n.inodes[2].key; string(k) != "foo" {
		t.Fatalf("exp=<foo>; got=<%s>", k)
	}
}