package tui

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	exec "example.com/downloader/exec"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type ProgressCallback func(downloaded int64)
type SetTotalCallback func(downloaded int64)
type SetDoneCallback func()

func DownloadFiles(fileNames []string, p *tea.Program) {
	const sshAlias = "bauer-prod-eu-cf-integration"
	fmt.Printf("Downloading from: %+v\n", sshAlias)

	for i, fileName := range fileNames {
		cwd, _ := os.Getwd()
		remote := "/client/dump/" + fileName
		local := filepath.Join(cwd, fileName)

		go DownloadFile(
			sshAlias,
			remote,
			local,
			func(total int64) {
				p.Send(setTotalMsg{Index: i, Total: total})
			},
			func(bytes int64) {
				p.Send(progressMsg{Index: i, Bytes: bytes})
			},
			func() {
				p.Send(doneMsg{Index: i})
			},
		)
	}
}

func DownloadFile(
	sshAlias,
	remotePath,
	localPath string,
	setTotal SetTotalCallback,
	progress ProgressCallback,
	setDone SetDoneCallback,
) error {
	h, err := exec.NewClientFromSshConfig(sshAlias)
	if err != nil {
		fmt.Printf("err: %+v\n", err) // output for debug
		return err
	}

	key, err := os.ReadFile(h.KeyPath)
	if err != nil {
		fmt.Printf("err: %+v\n", err) // output for debug
		return err
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return err
	}

	config := &ssh.ClientConfig{
		User:            h.User,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	conn, err := ssh.Dial("tcp", h.Hostname+":22", config)
	if err != nil {
		fmt.Printf("err: %+v\n", err) // output for debug
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
	setTotal(total)

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

	setDone()
	return nil
}
