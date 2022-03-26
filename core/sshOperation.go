package core

import (
	"fmt"
	"io"
	"os"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type Operation interface {
	RunShell() error
	Upload() error
}

type Shell struct {
	sshClient *ssh.Client
	cmd       string
	session   *ssh.Session
	err       error
}

type Upload struct {
	sftpClient *sftp.Client
	src        string
	dst        string
}

func NewShellRunner(sshClient *ssh.Client, cmd string) *Shell {
	return &Shell{
		sshClient: sshClient,
		cmd:       cmd,
		session:   &ssh.Session{},
		err:       nil,
	}
}

func NewUploader(sftpClient *sftp.Client, src string, dst string) *Upload {
	return &Upload{
		sftpClient: sftpClient,
		src:        src,
		dst:        dst,
	}
}

func (op *Shell) RunShell() (string, error) {
	session, err := op.sshClient.NewSession()
	if err != nil {
		return "", fmt.Errorf("创建远程shell失败：%v", err)
	}
	output, err := session.CombinedOutput(op.cmd)
	if err != nil {
		return "", fmt.Errorf("执行远程shell失败：%v", err)
	}

	return string(output), nil
}

func (op *Upload) Upload() error {
	srcValue, err := os.Open(op.src)
	if err != nil {
		return fmt.Errorf("发送到远程时，打开本地文件失败：%v", err)
	}
	dstValue, err := op.sftpClient.Create(op.dst)
	if err != nil {
		return fmt.Errorf("发送到远程时，创建远程文件失败：%v", err)
	}
	defer srcValue.Close()
	defer dstValue.Close()

	buf := make([]byte, 1024)
	for {
		n, err := srcValue.Read(buf)
		if err != nil {
			if err != io.EOF {
				return fmt.Errorf("读取本地文件错误：%v", err)
			} else {
				break
			}
		}
		_, err = dstValue.Write(buf[:n])
		if err != nil {
			return fmt.Errorf("写入远端文件错误：%v", err)
		}
	}

	return nil
}
