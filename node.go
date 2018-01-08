package HeraDB

import (
    "sort"
    "bytes"
    "strings"
    "fmt"
)

const (
    _t  = 2
)

type node struct {
    isLeaf      bool
    parent      *node
    children    nodes       // children is JYL
    key         []byte
    inodes      inodes      // children is
}

func NewNode (isLeaf bool) *node {
    return &node{isLeaf:isLeaf}
}

func (self *node) String() string {
    return self.root().Format(0)
}

func (self *node) Format(depth int) string {
    if self.isLeaf {
        var format string
        for _, inode := range self.inodes {
            format += strings.Repeat("  ", depth) + string(inode.key) + "\n"
        }
        return format
    }
    var format string
    for i, c := range self.children {
        if c == nil { continue}
        if i == len(self.children) - 1 {
            format += strings.Repeat("  ", depth) + "\n" + c.Format(depth+1)
        } else {
            format += strings.Repeat("  ", depth) + string(self.inodes[i].key) + "\n" + c.Format(depth+1)
        }
    }
    return format
}

func (self *node) root() *node {
    if self.parent == nil {
        return self
    }
    return self.parent.root()
}

func (self *node) Put(key []byte) {
    self.put(key)
}

// put inserts a key and value
func (self *node) put(key []byte) {
    root := self.root()
    if root.isFull() {
        newRoot := NewNode(false)
        root.parent = newRoot
        if root.parent == nil {
            fmt.Print("what")
        }
        newRoot.children = make(nodes, 2)
        newRoot.inodes = make(inodes, 1)
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
        self.inodes = append(self.inodes, inode{})
        copy(self.inodes[index+1:], self.inodes[index:])
        self.inodes[index].key = key
    } else {
        if self.children[index].isFull() {
            self.split(index)
            if bytes.Compare(self.inodes[index].key, key) == -1 { index++}
        }
        self.children[index].insert(key)
    }
}

func (self *node) split(index int) {
    nodey := self.children[index]
    nodez := NewNode(nodey.isLeaf)
    nodez.inodes = append(nodez.inodes, nodey.inodes[_t:]...)
    if !nodey.isLeaf {
        nodez.children = append(nodez.children, nodey.children[:_t+1]...)
    }
    self.children = append(self.children, &node{})
    copy(self.children[index+1:], self.children[index:])
    self.children[index+1] = nodez
    for x := len(self.inodes) - 1; x >= index; x-- {
        self.inodes[x + 1] = self.inodes[x]
    }
    self.inodes = append(self.inodes, inode{})
    self.inodes[index] = nodey.inodes[index]
    nodey.inodes = nodey.inodes[:_t]    // delete nodey element
    nodey.inodes = append(nodey.inodes[:index], nodey.inodes[index+1:]...)
}

func (self *node)isFull() bool{
    if len(self.inodes) == 2 * _t - 1 {
        return true
    }
    return false
}


type nodes []*node

type inode struct {
    key     []byte
    value   []byte
}
type inodes []inode