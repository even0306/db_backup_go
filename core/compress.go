package core

import (
	"bytes"
	"compress/gzip"
)

type compress interface {
	CompressFile(f []byte) error
}

type Gz struct {
	input bytes.Buffer
}

//保存备份的文件并压缩
func (file *Gz) CompressFile(f []byte, filepath string, filename string) {
	//压缩文件
	gf := gzip.NewWriter(&file.input)
	gf.Write(f)
	gf.Close()
}
