package main

import (
	"bufio"
	// "exec.go"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	// "golang.org/x/crypto/ssh"
	"strings"
)

func main() {
	host, err := NewClientFromSshConfig("bauer-prod-eu-cf-integration")
	if err != nil {
		panic(err)
	}

	prefix := "hockey_eu_"
	cmd := fmt.Sprintf(
		// `ls`,
		// `cd /client/EU/archive; ls`,
		`cd /client/EU/archive; ls | sed -n 's/%s[a-z_]*\.//p' | sed -n 's/\.csv//p' | sort | uniq | tail -n 100 | tac`,
		// `ls | sed -n 's/%s[a-z_]*\.//p' | sed -n 's/\.csv//p' | sort | uniq | tail -n 100 | tac`,
		prefix,
	)
	output, err := RunRemoteCommand(host, cmd)
	if err != nil {
		panic(err)
	}

	fmt.Println("REMOTE OUTPUT:")
	fmt.Println(output)

	lines := []string{}
	scanner := bufio.NewScanner(strings.NewReader(string(output)))

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}
	fmt.Printf("lines:\n")
	fmt.Printf("%+v\n", lines)

	return

	fileNames := []string{
		"hockey_eu_product.251103012539.csv",
		"hockey_eu_pricing.251103012539.csv",
		"hockey_eu_sku.251103012539.csv",
	}
	m := newModel(fileNames)
	p := tea.NewProgram(m)

	go DownloadFiles(fileNames, p)

	if _, err := p.Run(); err != nil {
		panic(err)
	}
	fmt.Printf("\nDownload complete! 🎉 \n")

	tea.Quit()
}
