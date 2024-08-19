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

type tree struct {
	Mode      string
	Path      string
	Sha       string
	ChildTree []*tree
}

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), err
}

func handleTreeNode(mode, path, sha string) *tree {
	return &tree{
		Mode:      mode,
		Path:      path,
		Sha:       sha,
		ChildTree: []*tree{},
	}
}

func (currTree *tree) serialise() ([]byte, error) {

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
	path := currtree.Path

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
		Mylog.Println(path + "/" + item.Name())
		if item.IsDir() {
			if strings.HasPrefix(item.Name(), ".git") {
				continue
			}
			childTree := handleTreeNode("tree", path+"/"+item.Name(), "")
			hashPathChildTree, err := childTree.writeObject(indexTable, save)
			if err != nil {
				Mylog.Println(err)
				return "", err
			}
			if len(childTree.ChildTree) == 0 {
				Mylog.Println("Passing here for : ", item.Name())
				continue
			}

			childTree.Sha = hashPathChildTree
			currtree.ChildTree = append(currtree.ChildTree, childTree)

			info, err := os.Stat(path + "/" + item.Name())
			if err != nil {
				Mylog.Println(err)
				continue
			}
			indexTable[path+"/"+item.Name()] = fileData{Path: hashPathChildTree, TimeModified: info.ModTime().String()}
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
			currtree.ChildTree = append(currtree.ChildTree, handleTreeNode("blob", path+"/"+item.Name(), sha))
		}
	}

	Mylog.Println(currtree)
	if err != nil {
		Mylog.Println(err)
		return "", err
	}
	encodedTree, err := currtree.serialise()

	Mylog.Println(path)
	for _, tr := range currtree.ChildTree {
		Mylog.Println(*tr)
	}

	Mylog.Println(string(encodedTree))

	//encryptedData := EncryptBytes(encodedTree)
	format := "tree"

	formatByte := []byte(format + " ")

	finalByteArray := append(formatByte, encodedTree...)
	encryptedData := EncryptBytes(finalByteArray)

	hashPath := ReturnHash(encryptedData.Bytes())
	//Mylog.Println("Writing ", string(finalByteArray), " at", hashPath, " at", path)

	if save {
		storeDataToFile(encryptedData, fileCreateOnly, true, get_repo(), ".gitbutworse", "objects", hashPath[:2], hashPath[2:])
	}
	currtree.Sha = hashPath

	return hashPath, nil
}
