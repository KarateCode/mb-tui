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
	peakEnv     peakEnv
}

func NewCopyBatchFilesModel(batchNumber string, env peakEnv) CopyBatchFilesModel {
	return CopyBatchFilesModel{
		batchNumber: batchNumber,
		peakEnv:     env,
	}
}

func (m CopyBatchFilesModel) Init() tea.Cmd {
	return func() tea.Msg {

		host, err := exec.NewClientFromSshConfig(m.peakEnv.sshServer)
		if err != nil {
			panic(err)
		}

		prefix := calcPrefix(m.peakEnv.clientCode)
		cmd := fmt.Sprintf(
			`cd /client/%s/archive; cp %s*%s.* /client/dump; cd /client/dump; ls %s*%s.*`,
			m.peakEnv.subFolder,
			prefix,
			m.batchNumber,
			prefix,
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
