package peak_setup_integration

import (
	// "fmt"
	// "strings"

	// "example.com/downloader/batchmenu"
	// "github.com/charmbracelet/bubbles/progress"
	batchmenu "example.com/downloader/batchmenu"
	downloader "example.com/downloader/tui/downloader"
	tea "github.com/charmbracelet/bubbletea"
)

type step int

const (
	stepBatchMenu step = iota
	stepDownloading
)

type Model struct {
	step       step
	batchMenu  batchmenu.Model
	downloader downloader.Model
	Program    *tea.Program
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func NewModel(lines []string) *Model {
	m := batchmenu.NewMenu(lines)
	return &Model{
		step:      0,
		batchMenu: m,
	}
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			// m.quitting = true
			return m, tea.Quit
		}
	}

	switch m.step {

	case stepBatchMenu:
		var cmd tea.Cmd
		// m.batchMenu, _ = m.batchMenu.Update(msg)
		m.batchMenu, cmd = m.batchMenu.Update(msg)

		// When batch is chosen -> transition
		if m.batchMenu.Done {
			fileNames := []string{
				"hockey_eu_product.251103012539.csv",
				"hockey_eu_pricing.251103012539.csv",
				"hockey_eu_sku.251103012539.csv",
			}
			m.downloader = downloader.NewModel(fileNames)
			// need to run DownloadFiles, but we're missing tea Program
			downloader.DownloadFiles(fileNames, m.Program)
			m.step = stepDownloading
			return m, m.downloader.Init()
		}

		return m, cmd

	case stepDownloading:
		var cmd tea.Cmd
		m.downloader, cmd = m.downloader.Update(msg)

		if m.downloader.Done {
			return m, tea.Quit
		}

		return m, cmd
	}

	return m, nil
}

func (m *Model) View() string {
	switch m.step {

	case stepBatchMenu:
		return m.batchMenu.View()

	case stepDownloading:
		return m.downloader.View()

	default:
		return "unknown state"
	}
}
