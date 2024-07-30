package main


import (
	"fmt"
	"os"
)



func main(){
    
    args := os.Args
    mydir , err := os.Getwd()
    if err != nil {
        fmt.Println(err)
        return
    }
    tmp := make(map[string][]byte)
    tmp["/Users/pratham/Desktop/cOdInG/gitFromScratch/src/gitFromScratch/go.mod"] = []byte{120,156,92,203,65,14,2,33,12,5,208,61,167,224,2,82,245,12,70,227,202,133,241,0,88,8,144,177,243,39,109,185,191,123,182,47,121,244,177,170,70,135,102,239,89,232,86,109,115,28,196,175,242,220,31,212,134,223,21,242,102,205,206,157,76,121,165,134,36,40,49,10,202,252,213,216,134,247,249,77,12,161,109,236,130,114,154,99,41,33,52,196,75,186,158,195,63,0,0,255,255,54,58,44,104} 
    cm1 := commit{
        id : "123",
        filePathAndContents: tmp,
    }

    rp := repo{

        basePath: mydir,
    }

    lcef := make(map[string]*commit , 0)
    lcef["/Users/pratham/Desktop/cOdInG/gitFromScratch/src/gitFromScratch/go.mod"] = &cm1
    b := branch{
        repo : rp,
        lastCommitEachFile: lcef,
    }
    
    switch args[1]{
        case "diff":
            b.get_diff()
    }
}


