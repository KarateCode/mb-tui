package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kevinburke/ssh_config"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type ProgressCallback func(downloaded int64)
type SetTotalCallback func(downloaded int64)

type ResolvedHost struct {
	Hostname string
	User     string
	Port     string
	KeyPath  string
}

func DownloadFiles(fileNames []string, p *tea.Program) {
	done := make([]chan bool, len(fileNames))
	for i := range done {
		done[i] = make(chan bool)
	}

	for i, fileName := range fileNames {
		ch := done[i]

		cwd, _ := os.Getwd()
		remote := "/client/dump/" + fileName
		local := filepath.Join(cwd, fileName)
		go DownloadFile(
			"bauer-prod-eu-cf-integration",
			remote,
			local,
			func(total int64) {
				p.Send(setTotalMsg(total))
			},
			func(bytes int64) {
				p.Send(progressMsg(bytes))
			},
			ch,
		)
	}

	fmt.Printf("attempting read from channels\n")

	for _, ch := range done {
		<-ch
		fmt.Printf("got one! read from a channel\n")
	}

}

func DownloadFile(
	sshAlias,
	remotePath,
	localPath string,
	setTotal SetTotalCallback,
	progress ProgressCallback,
	c chan bool,
) error {
	// Load ~/.ssh/config
	cfgPath := filepath.Join(os.Getenv("HOME"), ".ssh", "config")
	f, err := os.Open(cfgPath)
	if err != nil {
		return fmt.Errorf("ssh config open: %w", err)
	}
	defer f.Close()

	sshCfg, err := ssh_config.Decode(f)
	if err != nil {
		return fmt.Errorf("ssh config parse: %w", err)
	}

	fmt.Printf("Downloading from: %+v\n", sshAlias)

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

	fmt.Printf("Downloaded %d / %d bytes\n", downloaded, total)
	c <- true
	return nil
}
