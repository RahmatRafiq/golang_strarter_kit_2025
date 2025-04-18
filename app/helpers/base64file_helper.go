package helpers

import (
	"encoding/base64"
	"os"
)

func Base64FileToBytes(base64File string) ([]byte, error) {
	dec, err := base64.StdEncoding.DecodeString(base64File)
	if err != nil {
		return nil, err
	}

	return dec, nil
}

func StoreBase64File(base64File, path string, filename string) error {
	dec, err := Base64FileToBytes(base64File)
	if err != nil {
		return err
	}

	if len(path) == 0 {
		path = "/"
	}

	if len(filename) == 0 {
		filename = "file"
	}

	// add slash on suffix of path if not exists
	if path[len(path)-1:] != "/" {
		path += "/"
	}

	// remove slash on prefix of filename if exists
	if filename[:1] == "/" {
		filename = filename[1:]
	}

	err = os.MkdirAll(StoragePath()+path, os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(StoragePath() + path + filename)
	if err != nil {
		return err
	}

	defer file.Close()

	if _, err := file.Write(dec); err != nil {
		return err
	}

	if err := file.Sync(); err != nil {
		return err
	}

	return nil
}