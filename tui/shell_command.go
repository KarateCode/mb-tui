package tui

import (
	exec "example.com/downloader/exec"
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SshCmdCallback func(selected string) tea.Msg

type SshCmdModel struct {
	sshAlias   string
	sshCommand string
	message    string
	emitOutput SshCmdCallback
	spin       spinner.Model
}

func NewShellCmd(sshAlias string, sshCommand string, message string, callback SshCmdCallback) SshCmdModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205")) // magenta

	return SshCmdModel{
		sshAlias:   sshAlias,
		sshCommand: sshCommand,
		message:    message,
		emitOutput: callback,
		spin:       s,
	}
}

func (m SshCmdModel) Init() tea.Cmd {
	work := func() tea.Msg {
		host, err := exec.NewClientFromSshConfig(m.sshAlias)
		if err != nil {
			panic(err)
		}

		_, err = exec.RunRemoteCommand(host, m.sshCommand)
		output, err := exec.RunRemoteCommand(host, m.sshCommand)
		if err != nil {
			panic(err)
		}

		return m.emitOutput(output)
	}

	return tea.Batch(
		m.spin.Tick,
		work,
	)
}

func (m SshCmdModel) Update(msg tea.Msg) (SshCmdModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		}
	}

	m.spin, cmd = m.spin.Update(msg)
	return m, cmd
}

func (m SshCmdModel) View() string {
	return fmt.Sprintf(
		"\n\n  %s %s",
		m.spin.View(),
		m.message,
	)
}
