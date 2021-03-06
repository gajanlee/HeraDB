package HeraDB

import (
	"strings"
	"bytes"
	"sort"
)

type node struct {
	isLeaf		bool
	parent		*node
	children	[]*node
	inodes		[]inode
	keycount	int
}

type inode struct {
	key []byte
}

const (
	_t = 2
)

// appendINode must ensure node is not full
func (n *node) appendINode(in inode) {
	n.inodes[n.keycount] = in
	n.keycount++
}

func (n *node) mergeNode(nd *node) {
	for x := 0; x < nd.keycount; x++ {
		n.appendINode(nd.inodes[x])

		n.keycount++

	}
}

func (in *inode) reset() {
	in.key = nil
}

func (in *inode) isNIL() bool {
	return in.key == nil
}

func (n *node)Format(depth int) string {
	var format string
	if n.isLeaf {
		for _, inode := range n.inodes {
			if inode.key != nil {
				format += strings.Repeat("--", depth) + string(inode.key) + "\n"
			}
		}
	} else {
		if n.children[0] != nil {
			format = n.children[0].Format(depth+1)
		}
		for i, inode := range n.inodes {
			format += strings.Repeat("--", depth) + string(inode.key) + "\n"
			if n.children[i+1] != nil { format += n.children[i+1].Format(depth+1)}
		}
	}
	return format
}

func (n *node)String() string {
	return n.root().Format(0)
}

func (n *node)Put(key []byte) {
    n.root().preInsert(key)
}

func (n *node)Del(key []byte) {
    nd, i := n.Get(key)
    if nd != nil { nd.del(i)}
}

func newNode(isLeaf bool) *node {
	return &node{isLeaf: isLeaf, children: make([]*node, 2*_t), inodes: make([]inode, 2*_t-1)}
}

func CreateBTree() *node{
	return newNode(true)
}

func (n *node) root() *node{
	if n.parent == nil {
		return n
	}
	return n.parent.root()
}

func (n *node)del(index int) {
	// now this node should be deleted
	if n.isLeaf { n.inodes[index].reset(); n.keycount--; return }
	if n.children[index].keycount >= _t {
		preIndex := n.children[index].keycount - 1
		preNode := n.children[index].inodes[preIndex]

		n.inodes[index] = preNode				// copy 前驱结点 to current node
		n.children[index].del(preIndex)			// recursively delete
		// copy children pointer?
	} else if n.children[index+1].keycount >= _t {
		backNode := n.children[index+1].inodes[0]
		n.inodes[index] = backNode
		n.children[index+1].del(0)
	} else {
		key := n.inodes[index]
		copy(n.children[index:], n.children[index+1:])
		copy(n.inodes[index:], n.inodes[index+1:])
		n.keycount--
		nodey := n.children[index]
		nodez := n.children[index+1]

		i := nodey.keycount
		nodey.appendINode(key)
		for x := 0; x < nodez.keycount; x++ {
			nodey.children[nodey.keycount] = nodez.children[x]
			nodey.appendINode(nodez.inodes[x])
		}
		nodey.children[nodey.keycount] = nodez.children[nodez.keycount]
		// nodey.mergeNode(nodez)
		nodey.del(i)		// recursively delete key(index i)
	}



}

func (n *node)Get(key []byte) (*node, int) {
	return n.root().get(key)
}

func (n *node)get(key []byte) (*node, int) {
	index := sort.Search(n.keycount, func(i int) bool { return bytes.Compare(n.inodes[i].key, key) != -1})
	if index >= n.keycount && !n.isLeaf{
		return n.children[index].get(key)
	} else if index >= n.keycount && n.isLeaf {
		return nil, 0
	}
	switch bytes.Compare(n.inodes[index].key, key) {
	case 0: return n, index
	case -1: if n.isLeaf { return nil, 0}
		return n.children[index+1].get(key)
	case 1: if n.isLeaf { return nil, 0}
		return n.children[index].get(key)
	}
	return nil, 0
}

func (n *node)preInsert(key []byte) {
	root := n.root()
	if root.isFull() {
		newRoot := newNode(false)
		root.parent = newRoot
		newRoot.children[0] = root
		newRoot.split(0)
		newRoot.insert(key)
	} else {
		root.insert(key)
	}
}

func (n *node)insert(key []byte) {
	index := sort.Search(n.keycount, func(i int) bool { return bytes.Compare(n.inodes[i].key, key) != -1})
	if n.isLeaf {
        copy(n.inodes[index+1:], n.inodes[index:])
		n.inodes[index].key = key
		n.keycount++
	} else {
		if n.children[index].isFull() {
			n.split(index)
			if bytes.Compare(n.inodes[index].key, key) == -1 {
				index++
			}
		}
		n.children[index].insert(key)
	}
}

// node's ith child is full, but n is not full interior node.
func (n *node)split(i int) {
	copy(n.children[i+2:], n.children[i+1:])
	copy(n.inodes[i+1:], n.inodes[i:])
	n.children[i], n.children[i+1], n.inodes[i] = n.children[i].splitTwo()
	n.keycount++
}

func (n *node)splitTwo() (*node, *node, inode) {
	nodel := newNode(n.isLeaf)
	noder := newNode(n.isLeaf)

	copy(nodel.inodes[:_t-1], n.inodes[:_t-1]); nodel.keycount = _t-1
	copy(noder.inodes[:_t-1], n.inodes[_t:]); noder.keycount = _t-1
	if !n.isLeaf {
		copy(nodel.children[:_t], n.children[:_t])
		copy(noder.children[:_t], n.children[_t:])
	}
	return nodel, noder, n.inodes[_t-1]
}

func (n *node)isFull() bool {
	return n.keycount == 2*_t-1
}