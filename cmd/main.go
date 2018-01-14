package main


import (
	. ".."
	"fmt"
	"strconv"
)

func main() {
	n := CreateBTree()
	for x := 0; x < 10; x++ {
		n.Put([]byte(strconv.Itoa(x)))
	}
	//fmt.Printf("%s", n)
	n.Del([]byte("2"))
	fmt.Printf("%s", n)
}
