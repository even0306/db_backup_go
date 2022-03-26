package connection

import (
	"fmt"
	"mysql_backup_go/config"
	"net"
	"time"

	"golang.org/x/crypto/ssh"
	_ "golang.org/x/crypto/ssh"
)

type socket interface {
	Connect() (*ssh.Client, error)
}

type Ssh struct {
	connInfo *config.ConfigFile
	ip       string
	port     string
	user     string
	password string
}

func NewSshSocket(connInfo *config.ConfigFile) *Ssh {
	return &Ssh{
		connInfo: connInfo,
		ip:       connInfo.REMOTE_HOST,
		port:     connInfo.REMOTE_PORT,
		user:     connInfo.REMOTE_USER,
		password: connInfo.REMOTE_PASSWORD,
	}
}

func (sf *Ssh) Connect() (*ssh.Client, error) {
	sf.ip = sf.connInfo.REMOTE_HOST
	sf.port = sf.connInfo.REMOTE_PORT
	sf.user = sf.connInfo.REMOTE_USER
	sf.password = sf.connInfo.REMOTE_PASSWORD

	config := ssh.ClientConfig{
		User: sf.user,
		Auth: []ssh.AuthMethod{ssh.Password(sf.password)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout: 10 * time.Second,
	}

	sshClient, err := ssh.Dial("tcp", sf.ip, &config)
	if err != nil {
		return nil, fmt.Errorf("ssh连接失败：%v", err)
	}

	return sshClient, nil
}
