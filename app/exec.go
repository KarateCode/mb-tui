package main

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"os"
)

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
