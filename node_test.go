package HeraDB

import "testing"
import (
	"fmt"
	"strconv"
)

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

func TestGetNode(t *testing.T) {
	n := CreateBTree()
	n.Put([]byte("baz"))
	n.Put([]byte("foo"))
	n.Put([]byte("bar"))
	n.Put([]byte("foo"))
	if nd := n.Get([]byte("test")); nd != nil {
		t.Fatalf("exp=nil; got=<%s>", nd)
	}
	if nd := n.Get([]byte("baz")); nd == nil {
		t.Fatalf("exp=not nil; got=<%s>", nd)
	}
}

func TestDelNode(t *testing.T) {
	n := CreateBTree()
	for x := 0; x < 10; x++ {
		n.Put([]byte(strconv.Itoa(x)))
	}
	//fmt.Printf("%s", n)
	n.Del([]byte("2"))
	fmt.Printf("%s", n)
}