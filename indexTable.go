package main

import (
	"bytes"
	"encoding/gob"
	"os"
	"path/filepath"
)

type fileData struct {
	Path         string
	TimeModified string
}

func getCommitIndexTable(sha string) (map[string]fileData, error) {
	data, err := os.ReadFile(filepath.Join(get_repo(), ".gitbutworse", "ref", sha[:2], sha[2:]))
	if err != nil {
		return map[string]fileData{}, err
	}

	decodedData, err := decodeFile(*bytes.NewBuffer(data))
	if err != nil {
		return map[string]fileData{}, err
	}

	indexTable, err := deserialiseIndexTable(decodedData)
	if err != nil {
		return map[string]fileData{}, err
	}
	return indexTable, nil
}

func deserialiseIndexTable(data []byte) (map[string]fileData, error) {

	var b map[string]fileData
	dec := gob.NewDecoder(bytes.NewBuffer(data))

	err := dec.Decode(&b)
	if err != nil {
		Mylog.Fatal(err.Error())
		return b, err
	}

	return b, nil
}

func getLatestCommit() (string, error) {

	data, err := os.ReadFile(filepath.Join(get_repo(), ".gitbutworse", "ref", "HEAD"))
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func diffTreeWithCurrentState(tree1 string, indexTable1 map[string]fileData, path string) error {

	localTree, indexTable, err := getLocalTreeAndIndexTable(path)
	if err != nil {
		return err
	}

	if tree1 == localTree.Tree[0].Sha {
		Mylog.Println("Both the trees are equal. No Change")
	} else {

		// Two different trees with two different index table
		diffTwoIndexTable(indexTable1, indexTable)
	}
	return nil

}

func diffTwoIndexTable(indexTable1, indexTable2 map[string]fileData) {

	checked := make(map[string]bool)
	for fileName, fileData := range indexTable1 {
		fileData2, ok := indexTable2[fileName]

		if ok {
			checked[fileName] = true
			if fileData2.TimeModified == fileData.TimeModified {
				Mylog.Println("No change in file : ", fileName)
			} else {
				Mylog.Println("Change detected in file : ", fileName)
				//Mylog.Println(fileData.Path)

				//file1, err := readObject(fileData.Path)
				//if err != nil {
					//Mylog.Println("Error reading file in commit : ", fileName, " error : ", err)
				//}

				//file2, err := os.ReadFile(fileName)
				//if err != nil {
					//Mylog.Println("Error reading file in local repo : ", fileName, " error : ", err)
				//}

				//printDiffBytes(file1.data, file2)
			}
		} else {
			Mylog.Println("File deleted : ", fileName)
		}
	}

	for fileName := range indexTable2 {
		_, ok := checked[fileName]
		if !ok {
			Mylog.Println("File added : ", fileName)
		}
	}
}

func getLocalTreeAndIndexTable(path string) (tree, map[string]fileData, error) {

	currentLocalTree := tree{
		Tree: []treeNode{
			handleTreeNode("tree", path, ""),
		},
	}

	indexTree := make(map[string]fileData)

	hashPath, err := currentLocalTree.writeObject(indexTree, false)
	if err != nil {
		return tree{}, indexTree, err
	}

	currentLocalTree.Tree[0].Sha = hashPath

	return currentLocalTree, indexTree, nil
}

func serialiseIndexTable(indexTable map[string]fileData) ([]byte, error) {

	var b bytes.Buffer
	enc := gob.NewEncoder(&b)

	if err := enc.Encode(indexTable); err != nil {
		Mylog.Println("Error encoding struct : ", err)
		return []byte{}, err
	}

	return b.Bytes(), nil
}

func handleIndexTable(indexTable map[string]fileData) (bytes.Buffer, error) {

	data, err := serialiseIndexTable(indexTable)
	if err != nil {
		return bytes.Buffer{}, err
	}
	encryptedIndexTable := EncryptBytes(data)
	return encryptedIndexTable, nil
}
