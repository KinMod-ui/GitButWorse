package main

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var Mylog = log.New(os.Stderr, "GBW: ", log.LstdFlags|log.Lshortfile)

const (
	fileCreateOnly = iota
	fileWriteAppend
	fileWriteOverWrite
)

func readObject(sha string) (object, error) {
	var gitObject object

	Mylog.Println(sha)

	filepth := filepath.Join(get_repo(), ".gitbutworse", "objects", sha[:2], sha[2:])

	buf, err := os.ReadFile(filepth)

	if err != nil {
		Mylog.Println(err)
		return object{}, err
	}

	var decryptBuff []byte

	r, err := zlib.NewReader(bytes.NewBuffer(buf))

	if err != nil {
		Mylog.Println(err)
		return object{}, err
	}
	r.Close()

	decryptBuff, err = io.ReadAll(r)

	if err != nil {
		Mylog.Println(err)
		return object{}, err
	}

	idxOfFmt := bytes.Index(decryptBuff, []byte(" "))
	gitObject.format = string(decryptBuff[:idxOfFmt])
	gitObject.processData(decryptBuff[idxOfFmt+1:])

	return gitObject, nil
}

func (gObject *object) processData(data []byte) {
	switch gObject.format {
	case "blob":
		{
			gObject.data = append(gObject.data, data...)
		}
	case "tree":
		{
			Tree := deserialise(*bytes.NewBuffer(data))
			for _, subTree := range Tree.ChildTree {
				//if i == 0 {
				//continue
				//}
				object, err := readObject(subTree.Sha)
				if err != nil {
					Mylog.Println(err)
				}
				gObject.data = append(gObject.data, object.data...)
			}
		}
	}
}

func EncryptBytes(bts []byte) bytes.Buffer {
	var in bytes.Buffer
	w := zlib.NewWriter(&in)
	w.Write(bts)
	w.Close()

	return in
}

func decodeFile(file bytes.Buffer) ([]byte, error) {
	var out bytes.Buffer
	r, err := zlib.NewReader(&file)
	if err != nil {
		return []byte{}, err
	}
	r.Close()
	io.Copy(&out, r)
	//Mylog.Println("out" , out.Bytes())
	return out.Bytes(), nil
}

func ReturnHash(byteArray []byte) string {

	hasher := sha256.New()
	hasher.Write(byteArray)
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	return sha
}

func printDiffBytes(file1, file2 []byte) {
	//var buf1, buf2 bytes.Buffer
	//buf1 = *bytes.NewBuffer(file1)
	//buf2 = *bytes.NewBuffer(file2)

	//actualFileByte1, err := decodeFile(buf1)
	//if err != nil {
	//Mylog.Println(err)
	//return
	//}

	//actualFileByte2, err := decodeFile(buf2)
	//if err != nil {
	//Mylog.Println(err)
	//return
	//}

	actualFileString1 := string(file1)
	actualFileString2 := string(file2)

	scanner1 := bufio.NewScanner(strings.NewReader(actualFileString1))
	scanner2 := bufio.NewScanner(strings.NewReader(actualFileString2))
	Mylog.Println("+++ : file1, --- : file2")

	cnt := 0

	for scanner1.Scan() {
		cnt++
		if !scanner2.Scan() {
			Mylog.Println("Line:", cnt, " +++ ", scanner1.Text())
		} else {
			txt1 := scanner1.Text()
			txt2 := scanner2.Text()
			if txt1 != txt2 {
				Mylog.Println("Line:", cnt, " +++ ", scanner1.Text())
				Mylog.Println("Line:", cnt, " --- ", scanner2.Text())
			}
		}
	}

	for scanner2.Scan() {
		cnt++
		Mylog.Println("Line:", cnt, " --- ", scanner2.Text())
	}
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func storeDataToFile(data bytes.Buffer, appendMode int, pathCheck bool, path ...string) bool {
	if len(data.Bytes()) == 0 {
		Mylog.Println("No data to write")
		return false
	}
	if pathCheck && appendMode == fileCreateOnly {
		ret, err := exists(filepath.Join(strings.Join(path, "/")))
		if err != nil {
			Mylog.Println(err)
			return false
		}

		if ret {
			Mylog.Println("No change in file already in repo")
			return false
		} else {
			Mylog.Println("change in file already in repo")
		}
	}

	err := os.MkdirAll(filepath.Dir(strings.Join(path, "/")), os.ModePerm)
	if err != nil {
		Mylog.Println("recieved error", err)
		return false
	}

	var fileout *os.File
	if appendMode == fileWriteAppend {
		fileout, err = os.OpenFile(strings.Join(path, "/"), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	} else {
		fileout, err = os.OpenFile(strings.Join(path, "/"), os.O_CREATE|os.O_WRONLY, 0600)
	}

	if err != nil {
		Mylog.Println("recieved error", err)
		return false
	}

	defer func() {
		if err := fileout.Close(); err != nil {
			Mylog.Println("recieved error", err)
			return
		}
	}()

	buf := make([]byte, 1024)
	for {
		n, err := data.Read(buf)
		buf = buf[:n]
		if err != nil && err != io.EOF {
			panic(err)
		}

		if err == io.EOF {
			break
		}

		if _, err := fileout.Write(buf); err != nil {
			panic(err)
		}
	}

	return true
}

func get_repo() string {
	dir, err := os.Getwd()
	if err != nil {
		Mylog.Println("Error found here", err)
		return ""
	}
	return dir
}
