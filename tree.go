package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type treeNode struct {
	Mode string
	Path string
	Sha  string
}

type tree struct {
	Tree []treeNode
}

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), err
}

func handleTreeNode(mode, path, sha string) treeNode {
	return treeNode{
		Mode: mode,
		Path: path,
		Sha:  sha,
	}
}

func (currTree tree) serialise() ([]byte, error) {

	var b bytes.Buffer

	enc := gob.NewEncoder(&b)

	if err := enc.Encode(currTree); err != nil {
		Mylog.Println("Error encoding struct : ", err)
		return []byte{}, err
	}

	return b.Bytes(), nil
}

func deserialise(bts bytes.Buffer) tree {

	//Mylog.Println(string(bts.Bytes()))

	var b tree
	dec := gob.NewDecoder(&bts)

	err := dec.Decode(&b)
	if err != nil {
		Mylog.Fatal(err.Error())
		return b
	}

	//Mylog.Println(b)

	return b
}

func (currtree *tree) writeObject(indexTable map[string]fileData, save bool) (string, error) {
	path := currtree.Tree[0].Path

	items, err := os.ReadDir(path)

	if err != nil {
		Mylog.Println(err)
		return "", err
	}

	if len(items) == 0 {
		return "", nil
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Name() < items[j].Name()
	})

	for _, item := range items {
		if item.IsDir() {
			if strings.HasPrefix(item.Name(), ".git") {
				continue
			}
			childTree := tree{
				Tree: []treeNode{
					handleTreeNode("tree", path+"/"+item.Name(), ""),
				},
			}
			hashPathChildTree, err := childTree.writeObject(indexTable, save)
			if err != nil {
				Mylog.Println(err)
				return "", err
			}
			if len(childTree.Tree) == 1 {
				//Mylog.Println("Passing here for : ", item.Name())
				continue
			}
			childTree.Tree[0].Sha = hashPathChildTree
			currtree.Tree = append(currtree.Tree, childTree.Tree...)
		} else {
			fileDataByte, err := os.ReadFile(path + "/" + item.Name())
			if err != nil {
				fmt.Println(err, "here")
				continue
			}

			leafBlob := blob{
				object{
					format: "blob",
					data:   fileDataByte,
				},
			}

			sha, err := leafBlob.writeObject(filepath.Join(get_repo(), ".gitbutworse", "objects"), save)
			if err != nil {
				Mylog.Println(err)
				continue
			}
			info, err := os.Stat(path + "/" + item.Name())
			if err != nil {
				Mylog.Println(err)
				continue
			}
			indexTable[path+"/"+item.Name()] = fileData{Path: sha, TimeModified: info.ModTime().String()}
			currtree.Tree = append(currtree.Tree, handleTreeNode("blob", path+"/"+item.Name(), sha))
		}
	}

	encodedTree, err := currtree.serialise()
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	//encryptedData := EncryptBytes(encodedTree)
	format := "tree"

	formatByte := []byte(format + " ")

	finalByteArray := append(formatByte, encodedTree...)
	encryptedData := EncryptBytes(finalByteArray)

	hashPath := ReturnHash(encryptedData.Bytes())
	Mylog.Println("Writing ", string(finalByteArray), " at", hashPath, " at", path)

	if save {
		storeDataToFile(encryptedData, fileCreateOnly, true, get_repo(), ".gitbutworse", "objects", hashPath[:2], hashPath[2:])
	}
	currtree.Tree[0].Sha = hashPath

	return hashPath, nil
}
