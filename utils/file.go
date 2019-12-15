package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var File = NewFile()

func NewFile() *file {
	return &file{}
}

type file struct {
}

// file or path is exists
func (f *file) PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// read file data
func (f *file) ReadAll(path string) (data string, err error) {
	fi, err := os.Open(path)
	if err != nil {
		return
	}
	defer fi.Close()

	fd, err := ioutil.ReadAll(fi)
	return string(fd), nil
}

// write file
func (f *file) WriteFile(filename string, data string) (err error) {
	var dataByte = []byte(data)
	err = ioutil.WriteFile(filename, dataByte, 0666)
	if err != nil {
		return
	}
	return
}

// create file
func (f *file) CreateFile(filename string) error {
	newFile, err := os.Create(filename)
	defer newFile.Close()
	return err
}

// get dir all files
func (f *file) WalkDir(dirPth, suffix string) (files []string, err error) {
	files = make([]string, 0, 30)
	suffix = strings.ToUpper(suffix)
	err = filepath.Walk(dirPth, func(filename string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fi.IsDir() {
			return nil
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) {
			files = append(files, filename)
		}
		return nil
	})
	return files, err
}
