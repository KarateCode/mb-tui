package peak_setup_integration

import (
	// "fmt"
	// "strings"

	// "example.com/downloader/batchmenu"
	// "github.com/charmbracelet/bubbles/progress"
	// batchmenu "example.com/downloader/batchmenu"
	downloader "example.com/downloader/tui/downloader"
	tea "github.com/charmbracelet/bubbletea"
)

type step int

const (
	stepIntegrationMenu step = iota
	stepBatchMenu
	stepCopyingBatchFiles
	stepDownloading
)

type IntegrationMenuChoice string
type BatchChoice string
type copyCompleteMsg []string
type Model struct {
	step            step
	integrationMenu IntegrationMenuModel
	batchMenu       BatchModel
	copyBatchFiles  CopyBatchFilesModel
	downloader      downloader.Model
	Program         *tea.Program
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func NewModel() *Model {
	m := NewIntegrationMenu()
	return &Model{
		step:            0,
		integrationMenu: m,
	}
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case IntegrationMenuChoice:
		choice := IntegrationMenuChoice(msg)
		m.batchMenu = NewMenu(choice)
		m.step = stepBatchMenu
		return m, m.batchMenu.Init()

	case BatchChoice:
		choice := string(msg)
		m.copyBatchFiles = NewCopyBatchFilesModel(choice)
		m.step = stepCopyingBatchFiles
		return m, m.copyBatchFiles.Init()

	case copyCompleteMsg:
		fileNames := []string(msg)
		m.downloader = downloader.NewModel(fileNames)
		downloadFiles := func() tea.Msg {
			return downloader.DownloadFiles(fileNames, m.Program)
		}
		m.step = stepDownloading
		return m, downloadFiles

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		}
	}

	switch m.step {

	case stepIntegrationMenu:
		var cmd tea.Cmd
		m.integrationMenu, cmd = m.integrationMenu.Update(msg)
		return m, cmd

	case stepBatchMenu:
		var cmd tea.Cmd
		m.batchMenu, cmd = m.batchMenu.Update(msg)

		// When batch is chosen -> transition
		// if m.batchMenu.Done {
		// 	fileNames := []string{
		// 		"hockey_eu_product.251103012539.csv",
		// 		"hockey_eu_pricing.251103012539.csv",
		// 		"hockey_eu_sku.251103012539.csv",
		// 	}
		// 	m.downloader = downloader.NewModel(fileNames)
		// 	downloadFiles := func() tea.Msg {
		// 		return downloader.DownloadFiles(fileNames, m.Program)
		// 	}
		// 	m.step = stepDownloading
		// 	return m, downloadFiles
		// }

		return m, cmd

	case stepCopyingBatchFiles:
		var cmd tea.Cmd
		m.copyBatchFiles, cmd = m.copyBatchFiles.Update(msg)
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

	case stepIntegrationMenu:
		return m.integrationMenu.View()

	case stepBatchMenu:
		return m.batchMenu.View()

	case stepCopyingBatchFiles:
		return m.copyBatchFiles.View()

	case stepDownloading:
		return m.downloader.View()

	default:
		return "unknown state"
	}
}
