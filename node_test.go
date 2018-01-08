package HeraDB

import "testing"
import "fmt"

func TestCreateBTree(t *testing.T) {
	root := NewNode(false)
	root.put([]byte("lee"))
	root.put([]byte("lee1"))
	root.put([]byte("lee2"))
	root.put([]byte("lee3"))
	root.put([]byte("lee4"))
	root.put([]byte("lee5"))
	root.put([]byte("lee6"))
	root.put([]byte("lee7"))
	fmt.Printf("%s", root)
}