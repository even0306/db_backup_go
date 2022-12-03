package common

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"log"
)

type Compress interface {
	CompressFile() (*bytes.Buffer, error)
}

type mygzip struct {
	input    bytes.Buffer
	data     *[]byte
	filename *string
}

//初始化压缩保存功能，传入二进制文件流和文件名，返回*mygzip的结构体实例
func NewCompress(f *[]byte, filename *string) *mygzip {
	return &mygzip{
		data:     f,
		filename: filename,
	}
}

//压缩字节流，传入 *bytes.Buffer，返回error
func (file *mygzip) CompressFile() (*bytes.Buffer, error) {
	gwf := gzip.NewWriter(&file.input)
	gwf.Name = *file.filename
	_, err := gwf.Write(*file.data)
	if err != nil {
		return nil, fmt.Errorf("压缩数据失败：%w", err)
	}
	defer func() {
		err := gwf.Close()
		if err != nil {
			log.Panicf("压缩数据写入缓存关闭失败：%v", err)
		}
	}()

	return &file.input, nil
}
