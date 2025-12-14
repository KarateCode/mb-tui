package exec

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"bytes"
	"github.com/kevinburke/ssh_config"
	"golang.org/x/crypto/ssh"
)

type ResolvedHost struct {
	Hostname string
	User     string
	Port     string
	KeyPath  string
}

func NewClientFromSshConfig(sshAlias string) (*ResolvedHost, error) {
	// Load ~/.ssh/config
	cfgPath := filepath.Join(os.Getenv("HOME"), ".ssh", "config")
	f, err := os.Open(cfgPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	sshCfg, err := ssh_config.Decode(f)
	if err != nil {
		return nil, err
	}

	host := func(key string) string {
		v, _ := sshCfg.Get(sshAlias, key)
		return v
	}

	h := &ResolvedHost{
		Hostname: host("Hostname"),
		User:     host("User"),
		Port:     host("Port"),
		KeyPath:  host("IdentityFile"),
	}

	if h.Hostname == "" {
		h.Hostname = sshAlias // fallback, just like OpenSSH
	}
	if h.User == "" {
		h.User = "ec2-user" // optional default
	}
	if h.Port == "" {
		h.Port = "22"
	}
	if h.KeyPath != "" {
		h.KeyPath = strings.ReplaceAll(h.KeyPath, "~", "/Users/michaelschneider")
	}

	return h, nil
}

func RunRemoteCommand(host *ResolvedHost, command string) (string, error) {
	key, err := os.ReadFile(host.KeyPath)
	if err != nil {
		return "", fmt.Errorf("failed to read private key: %w", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %w", err)
	}

	config := &ssh.ClientConfig{
		User: host.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // (ok for internal networks)
	}

	conn, err := ssh.Dial("tcp", host.Hostname+":"+host.Port, config)
	if err != nil {
		return "", fmt.Errorf("failed to connect: %w", err)
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to open session: %w", err)
	}
	defer session.Close()

	var stdout bytes.Buffer
	session.Stdout = &stdout

	if err := session.Run(command); err != nil {
		return "", fmt.Errorf("command failed: %w", err)
	}

	return stdout.String(), nil
}

func RunLocalCommand(script string) {
	cmd := exec.Command("sh", "-c", script)
	out, err := cmd.CombinedOutput() // Captures both stdout and stderr
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	fmt.Printf("out:\n")
	fmt.Printf("%+v\n", out)
}
