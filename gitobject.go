package main

type object struct {
	data   []byte
	format string
}

type objectInter interface {
	writeObject() (string, error)
}
