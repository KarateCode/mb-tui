package scripts

import (
	"bufio"
	"fmt"
	"strings"

	"example.com/downloader/batchmenu"
	exec "example.com/downloader/exec"
	tea "github.com/charmbracelet/bubbletea"
)

func PeakSetupIntegration() {
	host, err := exec.NewClientFromSshConfig("bauer-prod-eu-cf-integration")
	if err != nil {
		panic(err)
	}

	prefix := "hockey_eu_"
	cmd := fmt.Sprintf(
		`cd /client/EU/archive; ls | sed -n 's/%s[a-z_]*\.//p' | sed -n 's/\.csv//p' | sort | uniq | tail -n 100 | tac`,
		prefix,
	)
	output, err := exec.RunRemoteCommand(host, cmd)
	if err != nil {
		panic(err)
	}

	lines := []string{}
	scanner := bufio.NewScanner(strings.NewReader(string(output)))

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}

	// lines := []string{"251203183900", "251203163900", "251203143900", "251203123900", "251203103900"}
	fmt.Printf("lines:\n")
	fmt.Printf("%+v\n", lines)

	m := batchmenu.NewMenu(lines)
	p := tea.NewProgram(m)
	if _, err := p.Run(); err == nil {
		// fmt.Println("Selected:", result.(batchmenu.Model).Selected())
		fmt.Println("exiting p.Run")
	}
}
