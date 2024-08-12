package main

import (
	"bytes"
	"compress/zlib"
	"crypto/sha256"
	"encoding/base64"
	"errors"
)

type blob struct {
	object
}

func (blob *blob) writeObject(path string, save bool) (string, error) {

	data := blob.data
	if len(data) == 0 {
		return "", errors.New("No data to write")
	}
	format := "blob"

	formatByte := []byte(format + " ")

	finalByteArray := append(formatByte, data[:]...)

	hasher := sha256.New()
	hasher.Write(finalByteArray)
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	var buff bytes.Buffer

	w := zlib.NewWriter(&buff)
	w.Write(finalByteArray)
	w.Close()

	if save {
		storeDataToFile(buff, true, path, sha[:2], sha[2:])
	}
	return sha, nil

}
