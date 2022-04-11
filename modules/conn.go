package modules

import (
	"fmt"
	"net"
	"time"

	"golang.org/x/crypto/ssh"
)

type Socket interface {
	Connect() (*ssh.Client, error)
}

type ConnInfo struct {
	Host     string
	Port     int
	User     string
	Password string
}

func NewSshSocket(host string, port int, user string, password string) *ConnInfo {
	return &ConnInfo{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
	}
}

func (sf *ConnInfo) Connect() (*ssh.Client, error) {
	config := ssh.ClientConfig{
		User: sf.User,
		Auth: []ssh.AuthMethod{ssh.Password(sf.Password)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout: 10 * time.Second,
	}

	sshClient, err := ssh.Dial("tcp", sf.Host+":"+string(sf.Port), &config)
	if err != nil {
		return nil, fmt.Errorf("ssh连接失败：%w", err)
	}

	return sshClient, nil
}
