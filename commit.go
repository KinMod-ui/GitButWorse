package main

import (
	"bytes"
	"encoding/gob"
)

type commit struct {
	Id           string
	IndexTable   map[string]fileData
	Tree         tree
	LastCommitId []string
	Message      string
}

func createCommitWithTreeAndIdxTable(sha string, tree tree, indexTable map[string]fileData, lastCommits []string,
	message string) error {

	newCommit := commit{
		Id:           sha,
		IndexTable:   indexTable,
		Tree:         tree,
		LastCommitId: lastCommits,
		Message:      message,
	}

	serialisedData, err := serialiseCommit(newCommit)
	if err != nil {
		Mylog.Println(err)
		return err
	}

	encryptedData := EncryptBytes(serialisedData.Bytes())

	currentLatestCommit, err := getLatestCommit()
	if err != nil {
		return err
	}

	storeDataToFile(*bytes.NewBuffer([]byte(sha)), fileWriteOverWrite, false, get_repo(), ".gitbutworse", "ref", "HEAD")

	ret := storeDataToFile(encryptedData, fileCreateOnly, true, get_repo(), ".gitbutworse", "ref", sha[:2], sha[2:])
	if ret == true {
		storeDataToFile(*bytes.NewBuffer([]byte("\n" + currentLatestCommit)), fileWriteAppend, false, get_repo(), ".gitbutworse",
			"refTable")
	}

	return nil
}

func deserialiseCommit(data []byte) (commit, error) {

	var b commit
	dec := gob.NewDecoder(bytes.NewBuffer(data))

	err := dec.Decode(&b)
	if err != nil {
		Mylog.Fatal(err.Error())
		return b, err
	}

	return b, nil
}

func serialiseCommit(cmt commit) (bytes.Buffer, error) {

	var b bytes.Buffer
	enc := gob.NewEncoder(&b)

	if err := enc.Encode(cmt); err != nil {
		Mylog.Println("Error encoding struct : ", err)
		return bytes.Buffer{}, err
	}

	return b, nil
}
