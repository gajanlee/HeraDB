package HeraDB

import (
    "sort"
    "bytes"
    "strings"
)

const (
    _t  = 2
)

type node struct {
    isLeaf      bool
    parent      *node
    children    nodes
    key         []byte
    inodes      inodes
}

func NewNode (isLeaf bool) *node {
    return &node{isLeaf:isLeaf}
}

func (self *node) String() string {
    return self.root().Format(0)
}

func (self *node) Format(depth int) string {
    if self.isLeaf {
        return "\n"
    }
    var format string
    for _, c := range self.children {
        format += strings.Repeat("  ", depth) + string(self.key) + "\n" + c.Format(depth+1)
    }
    return format
}

func (self *node) root() *node {
    if self.parent == nil {
        return self
    }
    return self.parent.root()
}

// put inserts a key and value
func (self *node) put(key []byte) {
    root := self.root()
    if root.isFull() {
        newRoot := NewNode(false)
        root.parent = newRoot
        newRoot.children[1] = root
        newRoot.split(1)
        newRoot.insert(key)
    } else {
        root.insert(key)
    }
}

func (self *node) insert(key []byte) {
    // Find insertion index.
    index := sort.Search(len(self.inodes), func(i int) bool { return bytes.Compare(self.inodes[i].key, key) != -1})
    if self.isLeaf {            
        self.inodes = append(self.inodes, &inode{})
        copy(self.inodes[index+1:], self.inodes[index:])
        copy(self.inodes[index].key, key)
    } else {
        if self.children[index].isFull() {
            self.split(index)
            if bytes.Compare(self.inodes[index].key, key) == -1 { index++}
        }
        self.children[index].insert(key)
    }
}

func (self *node) split(index int) {
    nodez := NewNode(false)
    nodey := self.children[index]
    nodez.isLeaf = nodey.isLeaf
    nodez.key = append(nodez.key, nodey.key[_t:]...)
    if !nodey.isLeaf {
        nodez.children = append(nodez.children, nodey.children[:_t+1]...)
    }
    copy(self.children[index+1:], self.children[index:])
    self.children[index+1] = nodez
    copy(self.key[index+1:], self.key[index:])
    self.key[index] = nodey.key[index]
}

func (self *node)isFull() bool{
    if len(self.children) == 2 * _t - 1 {
        return true
    }
    return false
}


type nodes []*node

type inode struct {
    key     []byte
    value   []byte
}
type inodes []*inode