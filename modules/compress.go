package modules

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

func NewCompress(f *[]byte, filename *string) *mygzip {
	return &mygzip{
		data:     f,
		filename: filename,
	}
}

//保存备份的文件并压缩
func (file *mygzip) CompressFile() (*bytes.Buffer, error) {
	//压缩文件
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
