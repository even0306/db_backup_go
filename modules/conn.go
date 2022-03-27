package modules

import (
	"fmt"
	"net"
	"time"

	"golang.org/x/crypto/ssh"
	_ "golang.org/x/crypto/ssh"
)

type Socket interface {
	Connect() (*ssh.Client, error)
}

type myssh struct {
	ip       string
	port     string
	user     string
	password string
}

func NewSshSocket(ip string, port string, user string, password string) *myssh {
	return &myssh{
		ip:       ip,
		port:     port,
		user:     user,
		password: password,
	}
}

func (sf *myssh) Connect() (*ssh.Client, error) {
	config := ssh.ClientConfig{
		User: sf.user,
		Auth: []ssh.AuthMethod{ssh.Password(sf.password)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout: 10 * time.Second,
	}

	sshClient, err := ssh.Dial("tcp", sf.ip+":"+sf.port, &config)
	if err != nil {
		return nil, fmt.Errorf("ssh连接失败：%w", err)
	}

	return sshClient, nil
}
