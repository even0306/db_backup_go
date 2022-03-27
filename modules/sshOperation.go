package modules

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/pkg/sftp"
)

type Operation interface {
	Upload() error
	Remove() error
}

type sftpInfo struct {
	sftpClient *sftp.Client
}

func NewSftpOperater(sftpClient *sftp.Client) *sftpInfo {
	return &sftpInfo{
		sftpClient: sftpClient,
	}
}

func (op *sftpInfo) Upload(src string, dst string) error {
	srcValue, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("发送到远程时，打开本地文件失败：%w", err)
	}
	dstValue, err := op.sftpClient.Create(dst)
	if err != nil {
		return fmt.Errorf("发送到远程时，创建远程文件失败：%w", err)
	}
	defer srcValue.Close()
	defer dstValue.Close()

	buf := make([]byte, 1024)
	for {
		n, err := srcValue.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) == false {
				return fmt.Errorf("读取本地文件错误：%w", err)
			} else {
				break
			}
		}
		_, err = dstValue.Write(buf[:n])
		if err != nil {
			return fmt.Errorf("写入远端文件错误：%w", err)
		}
	}

	return nil
}

func (op *sftpInfo) Remove(dst string) error {
	err := op.sftpClient.Remove(path.Join(dst))
	if err != nil {
		return fmt.Errorf("删除远程旧备份失败：%w", err)
	}
	return nil
}
