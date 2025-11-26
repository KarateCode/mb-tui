package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
	"github.com/pkg/sftp"
)

type ProgressCallback func(downloaded int64)

func DownloadFile(
	user, server, keyPath, remotePath, localPath string,
	progress ProgressCallback,
) error {
	key, err := os.ReadFile(keyPath)
	if err != nil {
		return err
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return err
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout: 10 * time.Second,
	}

	conn, err := ssh.Dial("tcp", server+":22", config)
	if err != nil {
		return err
	}
	defer conn.Close()

	sftpClient, err := sftp.NewClient(conn)
	if err != nil {
		return err
	}
	defer sftpClient.Close()

	src, err := sftpClient.Open(remotePath)
	if err != nil {
		return err
	}
	defer src.Close()

	stat, _ := src.Stat()
	total := stat.Size()

	dst, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	buf := make([]byte, 32*1024)
	var downloaded int64

	for {
		n, readErr := src.Read(buf)
		if n > 0 {
			_, writeErr := dst.Write(buf[:n])
			if writeErr != nil {
				return writeErr
			}
			downloaded += int64(n)
			progress(downloaded)
		}

		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			return readErr
		}
	}

	fmt.Printf("Downloaded %d / %d bytes\n", downloaded, total)
	return nil
}
