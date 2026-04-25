package conn

import (
	"fmt"
	"net"
	"time"

	"golang.org/x/crypto/ssh"
)

// 开始连接服务器，返回ssh客户端和error
func CreateSSHSocket(host string, port int, user string, password string) (*ssh.Client, error) {
	config := ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout: 10 * time.Second,
	}

	sshClient, err := ssh.Dial("tcp", host+":"+fmt.Sprint(port), &config)
	if err != nil {
		return nil, fmt.Errorf("ssh连接失败：%w", err)
	}

	return sshClient, nil
}
