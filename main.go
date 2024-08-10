package main

import (
	"fmt"
	"os"
)

func main() {

	var tp objectInter
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}

	tp = &tree{
		Tree: []treeNode{
			handleTreeNode("tree", dir, ""),
		},
	}

	hash, err := tp.writeObject()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(hash)
	t2, err := readObject(hash)
	fmt.Println(string(t2.data), t2.format)
}
