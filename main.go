package main

import (
	tea "github.com/charmbracelet/bubbletea"
	// "github.com/pkg/sftp"
)

func main() {
	remote := "/client/dump/hockey_eu_product.251103012539.csv"
	local := "/Users/michaelschneider/code/ti/somefile.csv"

	// First SSH/SFTP stat to get file size
	// (Skipping for brevity—you can add a function to stat the file)
	totalSize := int64(50_000_000) // pretend 50MB
	// stat, err := sftpClient.Stat(remotePath)
	// if err != nil {
	// 	return err
	// }
	// totalSize := stat.Size()

	m := newModel(totalSize)
	p := tea.NewProgram(m)

	go func() {
		DownloadFile(
			"ec2-user",
			"ec2-54-194-53-209.eu-west-1.compute.amazonaws.com",
			"/Users/michaelschneider/.ssh/bauer-prod-eu-cf.pem",
			remote,
			local,
			func(bytes int64) {
				p.Send(progressMsg(bytes))
			},
		)
	}()

	if _, err := p.Run(); err != nil {
		panic(err)
	}
}

// import (
// 	"fmt"
// 	"github.com/pkg/sftp"
// 	"golang.org/x/crypto/ssh"
// 	"io"
// 	"io/ioutil" // Use 'os' and 'io' for modern Go if you prefer
// 	"log"
// 	"os"
// ) // 088M

// func main() {
// 	// SFTP connection parameters
// 	host := "ec2-54-194-53-209.eu-west-1.compute.amazonaws.com"
// 	port := 22
// 	user := "ec2-user"
// 	// password := "your_password"
// 	remoteFile := "/client/dump/hockey_eu_product.251103012539.csv"
// 	localFile := "/Users/michaelschneider/code/ti/somefile.csv"
// 	pemFile := "/Users/michaelschneider/.ssh/bauer-prod-eu-cf.pem" // Path to your PEM file

// 	// Load the private key for authentication
// 	signer, err := getSigner(pemFile, "") // Pass the passphrase if needed
// 	if err != nil {
// 		log.Fatalf("Failed to get SSH signer from PEM file: %v", err)
// 	}

// 	// Create SSH client configuration using public key authentication
// 	config := &ssh.ClientConfig{
// 		User: user,
// 		Auth: []ssh.AuthMethod{
// 			ssh.PublicKeys(signer),
// 		},
// 		// InsecureIgnoreHostKey is used for simplicity;
// 		// for production, you should verify the host key.
// 		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
// 	}

// 	// Connect to the SSH server
// 	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, port), config)
// 	if err != nil {
// 		log.Fatalf("Failed to connect to SSH server: %v", err)
// 	}
// 	defer conn.Close()

// 	// Open SFTP session over the existing SSH connection
// 	client, err := sftp.NewClient(conn)
// 	if err != nil {
// 		log.Fatalf("Failed to open SFTP session: %v", err)
// 	}
// 	defer client.Close()

// 	// Download the file
// 	err = downloadFile(client, remoteFile, localFile)
// 	if err != nil {
// 		log.Fatalf("Failed to download file: %v", err)
// 	}

// 	fmt.Printf("Successfully downloaded [%s] to [%s]\n", remoteFile, localFile)
// }

// // getSigner reads a PEM file and returns an ssh.Signer
// func getSigner(file, passphrase string) (ssh.Signer, error) {
// 	// Read the key file content
// 	keyBytes, err := ioutil.ReadFile(file)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to read private key file: %w", err)
// 	}

// 	// Parse the key.
// 	// If a passphrase is required, use ssh.ParsePrivateKeyWithPassphrase(keyBytes, []byte(passphrase))
// 	signer, err := ssh.ParsePrivateKey(keyBytes)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to parse private key: %w", err)
// 	}

// 	return signer, nil
// }

// // downloadFile opens the remote file and copies its contents to a local file
// func downloadFile(client *sftp.Client, remoteFile, localFile string) (err error) {
// 	fmt.Printf("Downloading [%s] to [%s] ...\n", remoteFile, localFile)

// 	// Open the remote file
// 	srcFile, err := client.Open(remoteFile)
// 	if err != nil {
// 		return fmt.Errorf("unable to open remote file: %w", err)
// 	}
// 	defer srcFile.Close()

// 	// Create the local file
// 	dstFile, err := os.Create(localFile)
// 	if err != nil {
// 		return fmt.Errorf("unable to create local file: %w", err)
// 	}
// 	defer dstFile.Close()

// 	// Copy the file contents
// 	bytesCopied, err := io.Copy(dstFile, srcFile)
// 	if err != nil {
// 		return fmt.Errorf("unable to copy file contents: %w", err)
// 	}

// 	fmt.Printf("%d bytes copied\n", bytesCopied)
// 	return nil
// }
