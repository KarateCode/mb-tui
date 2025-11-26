package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	total     int64
	downloaded int64
	bar       progress.Model
	done      bool
}

type progressMsg int64

func newModel(total int64) model {
	return model{
		total: total,
		bar:   progress.New(progress.WithDefaultGradient()),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case progressMsg:
		m.downloaded = int64(msg)
		if m.downloaded >= m.total {
			m.done = true
		}
		return m, nil

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.done {
		return "Download complete! 🎉\n"
	}
	ratio := float64(m.downloaded) / float64(m.total)
	return fmt.Sprintf(
		"Downloading...\n\n%s\n\n%d%%\n",
		m.bar.ViewAs(ratio),
		int(ratio*100),
	)
}
