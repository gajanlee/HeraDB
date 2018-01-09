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
    n.root().Get(key).del(key)
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

func (n *node)del(key []byte) *node {
    return nil
}

func (n *node)Get(key []byte) *node{
    return n.root().get(key)
}

func (n *node)get(key []byte) *node {
    for i, inode := range n.inodes {
        switch bytes.Compare(key, inode.key) {
        case 1:
            if n.isLeaf {
                return nil
            } else {
                return n.children[i].get(key)
            }
        case 0: return n
        }
    }
    return nil
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