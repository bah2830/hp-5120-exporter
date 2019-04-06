package hpswitch

import (
	"fmt"
	"net"
	"strconv"

	"golang.org/x/crypto/ssh"

	"github.com/bah2830/hp-5120-exporter/pkg/networkswitch"
)

type Switch struct {
	clientConfig *ssh.ClientConfig
	client       *ssh.Client
	ip           string
	port         uint16
}

func NewWithPassword(ip string, port uint16, username, password string) (*Switch, error) {
	hpSwitch := &Switch{
		ip:   ip,
		port: port,
		clientConfig: &ssh.ClientConfig{
			User: username,
			Auth: []ssh.AuthMethod{
				ssh.Password(password),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Config: ssh.Config{
				Ciphers: []string{"aes128-cbc"},
			},
		},
	}

	client, err := ssh.Dial("tcp", net.JoinHostPort(ip, strconv.Itoa(int(port))), hpSwitch.clientConfig)
	if err != nil {
		return nil, err
	}

	hpSwitch.client = client

	return hpSwitch, nil
}

func (s *Switch) GetEnvironmentDetails() (networkswitch.EnvironmentDetails, error) {
	session, err := s.client.NewSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	output, err := session.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := session.Run("dis env\nquit\n"); err != nil {
		return nil, err
	}
	// output, err := session.CombinedOutput("dis env\n")
	// if err != nil {
	// 	return nil, err
	// }
	// fmt.Println(string(output))

	session.Wait()

	data := make([]byte, 0)
	if _, err := output.Read(data); err != nil {
		return nil, err
	}

	fmt.Println(string(data))
	return nil, nil
}

func (s *Switch) Close() {
	s.client.Close()
}
