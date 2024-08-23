## This is Git But Worse

Git but worse is a git implementation using Golang 

Commands implemented: 
```
git diff
git status
git commit
git checkout
git ls-tree
```

> "Everything in git is a git object" ~ Kin 2024

All the objects stored are first serialised, then encrypted and then stored in 
a file with the filepath of hash of the contents of serialised data

So in short:
```
Path = serialise(data) 
Contents = encrypted(serialised(data))
```

We have used 3 objects
1. Blob -> Storing file contents.
2. Tree -> Stores Array of tree nodes representing the current directory in git object form.
    The nodes if tree has a non-empty directory or otherwise represents a file.
    The nodes could either be made of blob or another Tree*
3. Commit -> Stores 4 things for now 
    a. Tree representing the current file structure at the time of commit 
    b. Index Table representing the filename mapped to a struct containing the last modified and 
        the virtual path storing the actual serialised encrypted contents of a file
    c. LastCommitId representing the id of last commit before this one
    d. Message representing the message of the commit

Blobs in the tree are sorted by the filePath so that no two paths can come in different order 
and create error results in git-diff
```
git diff -> Will diff the current tree and the last commited state.
            The files are traversed in the tree so only if hash of a node is changed we traverse that
            tree in a DFS-style manner.

```
```
git commit -> gets the latest tree and the indexTable and make a commit object based on all that data.
The index table and commit are stored using hash of the tree(Dont call me lazy. Call me efficient)
```

Some Common Paths
> ref/commitId -> stores a commit using tree hash

> ref/Head -> stores Head SHA

> ref/refTable -> stores all Hash inside a table for checkout commands

