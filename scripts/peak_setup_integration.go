package scripts

import (
	"bufio"
	// "exec.go"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	// "golang.org/x/crypto/ssh"
	// "batchmenu"
	// app "example.com/downloader/app"
	"example.com/downloader/batchmenu"
	exec "example.com/downloader/exec"

	"strings"
)

func PeakSetupIntegration() {
	host, err := exec.NewClientFromSshConfig("bauer-prod-eu-cf-integration")
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
	output, err := exec.RunRemoteCommand(host, cmd)
	if err != nil {
		panic(err)
	}

	// fmt.Println("REMOTE OUTPUT:")
	// fmt.Println(output)

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

	m := batchmenu.NewMenu(lines)
	p := tea.NewProgram(m)
	if _, err := p.Run(); err == nil {
		// fmt.Println("Selected:", result.(batchmenu.Model).Selected())
		fmt.Println("exiting p.Run")
	}
}
