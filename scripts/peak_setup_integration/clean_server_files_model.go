package peak_setup_integration

import (
	"fmt"

	exec "example.com/downloader/exec"
	tea "github.com/charmbracelet/bubbletea"
)

type CleanServerFilesModel struct {
	batchNumber string
	prefix      string
}

func NewCleanServerFilesModel(batchNumber string, prefix string) CleanServerFilesModel {
	return CleanServerFilesModel{
		batchNumber: batchNumber,
		prefix:      prefix,
	}
}

func (m CleanServerFilesModel) Init() tea.Cmd {
	return func() tea.Msg {
		host, err := exec.NewClientFromSshConfig("bauer-prod-eu-cf-integration")
		if err != nil {
			panic(err)
		}

		cmd := fmt.Sprintf(
			"cd /client/dump; rm %s*%s.*",
			m.prefix,
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
