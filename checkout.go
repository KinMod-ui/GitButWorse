package main

func checkoutBranch(hash string) error {

	commit, err := getCommitFromHash(hash)

	if err != nil {
		Mylog.Println(err)
		return err
	}

	//Mylog.Println(commit.IndexTable)

	lc, err := getLatestCommit()
	if err != nil {
		return err
	}
	//Mylog.Println(lc)

	latestCommit, err := getCommitFromHash(lc)
	if err != nil {
		return err
	}

	handleCommitDiff(commit, latestCommit)

	return nil
}

func handleCommitDiff(cmt1, cmt2 commit) {

	q := NewQueue()

	visited := make(map[string]bool)

	q.Push(&(cmt1.Tree))

	for {
		if q.len() == 0 {
			break
		}

		q2 := NewQueue()

		for {
			if q.len() == 0 {
				break
			}

			elem := (q.Pop()).(*tree)
			Mylog.Println(elem)

			hashInCmt2, ok := cmt2.IndexTable[elem.Path]

			if !ok {
				Mylog.Println("New file added : ", elem.Path)
				AddPath(elem)
				//for _, newPath := range elem.ChildTree {
				////Mylog.Println(newPath.serialise())
				//q2.Push(newPath)
				//}
			} else {
				hashInCmt1, _ := cmt1.IndexTable[elem.Path]
				if hashInCmt1 != hashInCmt2 {
					Mylog.Println("File changed at :", elem.Path)
					AddPath(elem)
					for _, newPath := range elem.ChildTree {
						//Mylog.Println(newPath.serialise())
						q2.Push(newPath)
					}
				} else {
					Mylog.Println("File is same at : ", elem.Path)
				}
			}
			visited[elem.Path] = true
		}
		q = q2
	}
}

func AddPath(currTree *tree) error {

	if currTree.Mode == "blob" {
		obj, err := readObject(currTree.Sha)
		if err != nil {
			return err
		}
		Mylog.Println(obj)
	}
	return nil
}
