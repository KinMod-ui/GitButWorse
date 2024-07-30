package main

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func getFiles(path string) []string{
    DirList := make([]string, 0)
    DirList = append(DirList, path)
    FileList := make([]string, 0)
    for len(DirList) > 0 {
        currDir := DirList[0]
        DirList = DirList[1:]
        files , err := os.ReadDir(currDir)
        if filepath.Base(currDir) == ".git"{
            continue
        }
        if err != nil{
            fmt.Println(err)
            return []string{}
        }
        for _ , file := range(files){
            if file.IsDir(){
                DirList = append(DirList, filepath.Join(currDir , file.Name()))
            }else{
                FileList = append(FileList, filepath.Join(currDir , file.Name()))
            }
        }
    }
    return FileList
}

func encodeFile(path string)(bytes.Buffer , error){

    file , err := os.ReadFile(path)
    if err != nil {
        fmt.Println(err)
        return bytes.Buffer{} , err
    } 
    path += "  "

    var buff bytes.Buffer

    pathByteArrayWithSpace:= []byte(path)
    file = append(pathByteArrayWithSpace[:] , file[:]...)

    w := zlib.NewWriter(&buff)
    w.Write(file)
    w.Close()
    return buff , nil
}


func decodeFile(file bytes.Buffer) ([]byte,error){
        var out bytes.Buffer
        r , err := zlib.NewReader(&file)
        if err != nil{
            return []byte{}, err
        }
        defer r.Close()
            io.Copy(&out , r)
        //fmt.Println("out" , out.Bytes())
    return out.Bytes() , nil
}


func printDiffBytes(file1 , file2 []byte){
    var buf1 , buf2 bytes.Buffer
    buf1 = *bytes.NewBuffer(file1) 
    buf2 = *bytes.NewBuffer(file2)

    actualFileByte1 , err := decodeFile(buf1)
    if err != nil {
        fmt.Println(err)
        return
    }

    actualFileByte2 , err := decodeFile(buf2)
    if err != nil {
        fmt.Println(err)
        return
    }

    actualFileString1 := string(actualFileByte1)    
    actualFileString2 := string(actualFileByte2)

    scanner1 := bufio.NewScanner(strings.NewReader(actualFileString1))
    scanner2 := bufio.NewScanner(strings.NewReader(actualFileString2))
    fmt.Println("+++ : file1, --- : file2")

    cnt := 0

    for scanner1.Scan(){
        cnt++
        if (!scanner2.Scan()){
            fmt.Println("Line:", cnt , " +++ " , scanner1.Text())
        }else {
            txt1 := scanner1.Text()
            txt2 := scanner2.Text()
            if (txt1 != txt2){
                fmt.Println("Line:", cnt , " +++ " , scanner1.Text())
                fmt.Println("Line:" ,cnt , " --- " , scanner2.Text())
            }
        }
    }

    for scanner2.Scan(){
        fmt.Println("Line:" ,cnt , " --- " , scanner2.Text())
    }
} 
