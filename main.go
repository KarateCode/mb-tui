package main

import (
	// "exec.go"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	// "golang.org/x/crypto/ssh"
)

func main() {
	host, err := NewClientFromSshConfig("bauer-prod-eu-cf-integration")
	if err != nil {
		panic(err)
	}

	// cmd := `ls | sed -n 's/${prefix}[a-z_]*\.//p' | sed -n 's/\.csv//p' | sort | uniq | tail -n 100 | tac`
	prefix := "hockey_eu_"
	cmd := fmt.Sprintf(
		`ls | sed -n 's/%s[a-z_]*\.//p' | sed -n 's/\.csv//p' | sort | uniq | tail -n 100 | tac`,
		prefix,
	)
	output, err := RunRemoteCommand(host, "~/.ssh/id_rsa", cmd)
	if err != nil {
		panic(err)
	}

	fmt.Println("REMOTE OUTPUT:")
	fmt.Println(output)

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
