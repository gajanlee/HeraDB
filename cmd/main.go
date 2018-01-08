package main


import (
	. ".."
	"fmt"
)

func main() {
	root := NewNode(true)
	root.Put([]byte("lee"))
	root.Put([]byte("lee1"))
	root.Put([]byte("lee2"))
	root.Put([]byte("lee3"))
	root.Put([]byte("lee4"))
	root.Put([]byte("lee5"))

	fmt.Printf("%s", root)
}
