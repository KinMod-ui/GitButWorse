package main

import (
	"os"
)

func main() {

	var tp *tree

	dir, err := os.Getwd()
	if err != nil {
		Mylog.Println(err)
		return
	}
	tp = &tree{
		Tree: []treeNode{
			handleTreeNode("tree", dir, ""),
		},
	}

	var indexTable = make(map[string]fileData)
	hash, err := tp.writeObject(indexTable, true)
	if err != nil {
		Mylog.Println(err)
		return
	}

	createCommitWithTreeAndIdxTable(hash, *tp, indexTable, []string{}, "Hi I am kin")

	getCommitFromHash(hash)
	//storeDataToFile(*bytes.NewBuffer([]byte(hash)), false, get_repo(), ".gitbutworse", "ref", "HEAD")

	//encryptedIndexTable, err := handleIndexTable(indexTable)
	//if err != nil {
	//Mylog.Println(err)
	//return
	//}

	//storeDataToFile(encryptedIndexTable, true, get_repo(), ".gitbutworse", "ref", hash[:2], hash[2:])

	//args := os.Args[1:]

	//switch args[0] {
	//case "diff":
	//{
	//currentHead, err := getLatestCommit()
	//if err != nil {
	//if os.IsNotExist(err) {
	//Mylog.Println(err)
	//} else {
	//Mylog.Println("Shouldnt have happened: ", err)
	//}
	//} else {
	//dir, err := os.Getwd()
	//if err != nil {
	//Mylog.Println(err)
	//return
	//}
	//headIndexTable, err := getCommitIndexTable(currentHead)
	//if err != nil {
	//Mylog.Println(err)
	//return
	//}

	//diffTreeWithCurrentState(currentHead, headIndexTable, dir)
	//}
	//}
	//}
}
