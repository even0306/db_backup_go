package core

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"log"
)

type compress interface {
	CompressFile() (*bytes.Buffer, error)
}

type Gzip struct {
	input    bytes.Buffer
	data     *[]byte
	filename *string
}

func NewCompress(f *[]byte, filename *string) *Gzip {
	return &Gzip{
		data:     f,
		filename: filename,
	}
}

//保存备份的文件并压缩
func (file *Gzip) CompressFile() (*bytes.Buffer, error) {
	//压缩文件
	gwf := gzip.NewWriter(&file.input)
	gwf.Name = *file.filename
	_, err := gwf.Write(*file.data)
	if err != nil {
		return nil, fmt.Errorf("压缩数据失败：%v", err)
	}
	defer func() {
		err := gwf.Close()
		if err != nil {
			log.Panicf("压缩数据写入缓存关闭失败：%v", err)
		}
	}()

	return &file.input, nil
}
