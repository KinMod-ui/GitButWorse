package main

//import (
	//"bytes"
	//"fmt"
	//"io"
	//"os"
//)

//type branch struct{
    //lastCommitEachFile map[string]*commit
    //commits []*commit
    //repo repo
//}

//func (b branch)get_diff(path ...string){
    //var root string
    //if len(path) > 0 {
        //root = path[0]
    //} else {
        //root = b.repo.get_root()
    //}

    //files := getFiles(root)
    
    //for _ , file := range(files){
        //commit, ok := b.lastCommitEachFile[file]
        //if (!ok){
            //encodeFile(file)
            //fmt.Println(file)
            //fileHandler , err := os.Open(file) 
            //if err != nil {
                //fmt.Printf("%s" , err.Error())
                //defer func(){
                    //if err = fileHandler.Close(); err != nil {
                        //fmt.Printf("%s" , err)
                    //}
                //}()
            //}

            //buf := make([]byte , 32*1024)
            
            //for {
                //n , err := fileHandler.Read(buf)

                //if err == io.EOF{
                    //break

                //}

                //if err != nil {
                    //fmt.Printf("read %d bytes: %s" , n , err)
                    //return
                //}
            //}
        //} else {
            //originalHash , ok := commit.filePathAndContents[file]
            //if !ok{
                //fmt.Println("Woah not supposed to be here")
                //return
            //} else {
                //newHash ,err := encodeFile(file)
                //if err != nil {
                    //fmt.Println(err)
                    //return
                //}
                //fmt.Println(file)
                //if (!bytes.Equal(newHash.Bytes() , originalHash)){
                   //printDiffBytes(newHash.Bytes() , originalHash) 
                //}
            //}
        //}
    //}
//}
