package peak_setup_integration

import (
	"fmt"

	exec "example.com/downloader/exec"
	tea "github.com/charmbracelet/bubbletea"
)

type CleanServerFilesModel struct {
	batchNumber string
	peakEnv     peakEnv
}

func NewCleanServerFilesModel(batchNumber string, env peakEnv) CleanServerFilesModel {
	return CleanServerFilesModel{
		batchNumber: batchNumber,
		peakEnv:     env,
	}
}

func (m CleanServerFilesModel) Init() tea.Cmd {
	return func() tea.Msg {
		host, err := exec.NewClientFromSshConfig(m.peakEnv.sshServer)
		if err != nil {
			panic(err)
		}

		prefix := calcPrefix(m.peakEnv.clientCode)
		cmd := fmt.Sprintf(
			"cd /client/dump; rm %s*%s.*",
			prefix,
			m.batchNumber,
		)

		output, err := exec.RunRemoteCommand(host, cmd)
		if err != nil {
			panic(err)
		}

		return cleanServerCompleteMsg(output)
	}
}

func (m CleanServerFilesModel) Update(msg tea.Msg) (CleanServerFilesModel, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m CleanServerFilesModel) View() string {
	return "\n\nCleaning up files on the server..."
}
