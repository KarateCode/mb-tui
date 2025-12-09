package tui

import (
	exec "example.com/downloader/exec"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
)

type SshCmdModel struct {
	sshAlias   string
	sshCommand string
	message    string
	emitOutput TeaCmdCallback
}

func NewShellCmd(sshAlias string, sshCommand string, message string, callback TeaCmdCallback) SshCmdModel {
	return SshCmdModel{
		sshAlias:   sshAlias,
		sshCommand: sshCommand,
		message:    message,
		emitOutput: callback,
	}
}

func (m SshCmdModel) Init() tea.Cmd {
	host, err := exec.NewClientFromSshConfig(m.sshAlias)
	if err != nil {
		panic(err)
	}

	output, err := exec.RunRemoteCommand(host, m.sshCommand)
	if err != nil {
		panic(err)
	}

	return m.emitOutput(output)
}

func (m SshCmdModel) Update(msg tea.Msg) (SshCmdModel, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m SshCmdModel) View() string {
	return fmt.Sprintf("\n\n%s", m.message)
}
