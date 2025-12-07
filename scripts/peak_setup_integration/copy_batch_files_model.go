package peak_setup_integration

import (
	"bufio"
	"fmt"
	"strings"

	exec "example.com/downloader/exec"
	tea "github.com/charmbracelet/bubbletea"
)

type CopyBatchFilesModel struct {
	batchNumber string
	prefix      string
}

func NewCopyBatchFilesModel(batchNumber string, prefix string) CopyBatchFilesModel {
	return CopyBatchFilesModel{
		batchNumber: batchNumber,
		prefix:      prefix,
	}
}

func (m CopyBatchFilesModel) Init() tea.Cmd {
	return func() tea.Msg {

		host, err := exec.NewClientFromSshConfig("bauer-prod-eu-cf-integration")
		if err != nil {
			panic(err)
		}

		cmd := fmt.Sprintf(
			`cd /client/EU/archive; cp %s*%s.* /client/dump; cd /client/dump; ls %s*%s.*`,
			m.prefix,
			m.batchNumber,
			m.prefix,
			m.batchNumber,
		)

		output, err := exec.RunRemoteCommand(host, cmd)
		if err != nil {
			panic(err)
		}

		lines := copyCompleteMsg{}
		scanner := bufio.NewScanner(strings.NewReader(string(output)))
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line != "" {
				lines = append(lines, line)
			}
		}

		return lines
	}
}

func (m CopyBatchFilesModel) Update(msg tea.Msg) (CopyBatchFilesModel, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m CopyBatchFilesModel) View() string {
	return "\n\nCopying batch files to /client/dumps..."
}
