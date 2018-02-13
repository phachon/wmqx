package utils

import (
	"os"
	"io/ioutil"
)

func NewFile() *File {
	return &File{}
}

type File struct {
	
}

// file or path is exists
func (f *File) PathExists(path string) (bool, error) {
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
func (f *File) ReadAll(path string) (data string, err error) {
	fi, err := os.Open(path)
	if err != nil {
		return
	}
	defer fi.Close()

	fd, err := ioutil.ReadAll(fi)
	return string(fd), nil
}


// write file
func (f *File) WriteFile(filename string, data string) (err error) {
	var dataByte = []byte(data)
	err = ioutil.WriteFile(filename, dataByte, 0666)
	if err != nil {
		return
	}
	return
}

// create file
func (f *File) CreateFile(filename string) error {
	newFile, err := os.Create(filename)
	defer newFile.Close()
	return err
}