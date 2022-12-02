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

//初始化连接器，传入远端服务器主机ip，端口，用户名，密码，返回*ConnInfo的结构体信息指针
func NewSshSocket(host string, port int, user string, password string) *ConnInfo {
	return &ConnInfo{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
	}
}

//开始连接服务器，返回ssh客户端和error
func (sf *ConnInfo) Connect() (*ssh.Client, error) {
	config := ssh.ClientConfig{
		User: sf.User,
		Auth: []ssh.AuthMethod{ssh.Password(sf.Password)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout: 10 * time.Second,
	}

	sshClient, err := ssh.Dial("tcp", sf.Host+":"+fmt.Sprint(sf.Port), &config)
	if err != nil {
		return nil, fmt.Errorf("ssh连接失败：%w", err)
	}

	return sshClient, nil
}
