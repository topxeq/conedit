package editor

import (
	"fmt"
	"io/ioutil"
	"strings"

	"golang.org/x/crypto/ssh"
)

type SSHConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	KeyPath  string
}

type SSHClient struct {
	config *SSHConfig
	client *ssh.Client
}

func NewSSHClient(config *SSHConfig) *SSHClient {
	return &SSHClient{config: config}
}

func (s *SSHClient) Connect() error {
	authMethods := []ssh.AuthMethod{}

	if s.config.Password != "" {
		authMethods = append(authMethods, ssh.Password(s.config.Password))
	}

	if s.config.KeyPath != "" {
		key, err := ioutil.ReadFile(s.config.KeyPath)
		if err == nil {
			signer, err := ssh.ParsePrivateKey(key)
			if err == nil {
				authMethods = append(authMethods, ssh.PublicKeys(signer))
			}
		}
	}

	if len(authMethods) == 0 {
		return fmt.Errorf("no authentication method available")
	}

	clientConfig := &ssh.ClientConfig{
		User:            s.config.User,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	addr := s.config.Host
	if s.config.Port != "" {
		addr = addr + ":" + s.config.Port
	} else {
		addr = addr + ":22"
	}

	client, err := ssh.Dial("tcp", addr, clientConfig)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	s.client = client
	return nil
}

func (s *SSHClient) ReadFile(path string) (string, error) {
	if s.client == nil {
		return "", fmt.Errorf("not connected")
	}

	session, err := s.client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	output, err := session.CombinedOutput("cat " + path)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	return string(output), nil
}

func (s *SSHClient) WriteFile(path, content string) error {
	if s.client == nil {
		return fmt.Errorf("not connected")
	}

	session, err := s.client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	session.Stdin = strings.NewReader(content)
	err = session.Run("cat > " + path)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func (s *SSHClient) Close() {
	if s.client != nil {
		s.client.Close()
	}
}
