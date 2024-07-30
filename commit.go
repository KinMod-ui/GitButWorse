package main


type commit struct{
    id string
    filePathAndContents map[string][]byte
    lastCommitId map[string]*commit
    message string
}
