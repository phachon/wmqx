package utils

import (
	"bytes"
	"os"
)

const (
	defaultBufSize = 4096
)

var Tail = NewTail()

type tail struct{}

func NewTail() *tail {
	return &tail{}
}

func (t *tail) Run(filename string, n int) (lines []string, err error) {
	f, e := os.Stat(filename)
	if e == nil {
		size := f.Size()
		var fi *os.File
		fi, err = os.Open(filename)
		if err == nil {
			b := make([]byte, defaultBufSize)
			sz := int64(defaultBufSize)
			nn := n
			bTail := bytes.NewBuffer([]byte{})
			istart := size
			flag := true
			for flag {
				if istart < defaultBufSize {
					sz = istart
					istart = 0
					//flag = false
				} else {
					istart -= sz
				}
				_, err = fi.Seek(istart, os.SEEK_SET)
				if err == nil {
					mm, e := fi.Read(b)
					if e == nil && mm > 0 {
						j := mm
						for i := mm - 1; i >= 0; i-- {
							if b[i] == '\n' {
								bLine := bytes.NewBuffer([]byte{})
								bLine.Write(b[i+1 : j])
								j = i
								if bTail.Len() > 0 {
									bLine.Write(bTail.Bytes())
									bTail.Reset()
								}

								if (nn == n && bLine.Len() > 0) || nn < n {
									//skip last "\n"
									lines = append(lines, bLine.String())
									nn--
								}
								if nn == 0 {
									flag = false
									break
								}
							}
						}
						if flag && j > 0 {
							if istart == 0 {
								bLine := bytes.NewBuffer([]byte{})
								bLine.Write(b[:j])
								if bTail.Len() > 0 {
									bLine.Write(bTail.Bytes())
									bTail.Reset()
								}
								lines = append(lines, bLine.String())
								flag = false
							} else {
								bb := make([]byte, bTail.Len())
								copy(bb, bTail.Bytes())
								bTail.Reset()
								bTail.Write(b[:j])
								bTail.Write(bb)
							}
						}
					}
				}
			}
			//func (f *File) Seek(offset int64, whence int) (ret int64, err error)
			//func (f *File) Read(b []byte) (n int, err error) {
		}
		defer fi.Close()
	}
	return
}
