package peak_setup_integration

import (
	// "fmt"
	// "strings"

	// "example.com/downloader/batchmenu"
	// "github.com/charmbracelet/bubbles/progress"
	// batchmenu "example.com/downloader/batchmenu"
	"fmt"

	downloader "example.com/downloader/tui/downloader"
	tea "github.com/charmbracelet/bubbletea"
)

type step int

const (
	stepEnvMenu step = iota
	stepIntegrationMenu
	stepBatchMenu
	stepCopyingBatchFiles
	stepDownloading
	stepCleanServerFiles
)

type IntegrationMenuChoice string
type BatchChoice string
type EnvMenuChoice peakEnv
type copyCompleteMsg []string
type cleanServerCompleteMsg string
type Model struct {
	step             step
	peakEnvMenu      EnvMenuModel
	integrationMenu  IntegrationMenuModel
	batchMenu        BatchModel
	copyBatchFiles   CopyBatchFilesModel
	downloader       downloader.Model
	cleanServerFiles CleanServerFilesModel

	batchChoice   string
	envMenuChoice peakEnv
	Program       *tea.Program
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func NewModel() *Model {
	m := NewEnvMenu()
	return &Model{
		step:        0,
		peakEnvMenu: m,
	}
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case EnvMenuChoice:
		choice := peakEnv(msg)
		m.envMenuChoice = choice
		m.integrationMenu = NewIntegrationMenu()
		m.step = stepIntegrationMenu
		return m, m.integrationMenu.Init()

	case IntegrationMenuChoice:
		choice := IntegrationMenuChoice(msg)
		m.batchMenu = NewMenu(choice)
		m.step = stepBatchMenu
		return m, m.batchMenu.Init()

	case BatchChoice:
		choice := string(msg)
		m.batchChoice = choice
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

	case cleanServerCompleteMsg:
		fmt.Print("Download (app) complete!")
		return m, tea.Quit

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		}
	}

	switch m.step {

	case stepEnvMenu:
		var cmd tea.Cmd
		m.peakEnvMenu, cmd = m.peakEnvMenu.Update(msg)
		return m, cmd

	case stepIntegrationMenu:
		var cmd tea.Cmd
		m.integrationMenu, cmd = m.integrationMenu.Update(msg)
		return m, cmd

	case stepBatchMenu:
		var cmd tea.Cmd
		m.batchMenu, cmd = m.batchMenu.Update(msg)
		return m, cmd

	case stepCopyingBatchFiles:
		var cmd tea.Cmd
		m.copyBatchFiles, cmd = m.copyBatchFiles.Update(msg)
		return m, cmd

	case stepDownloading:
		var cmd tea.Cmd
		m.downloader, cmd = m.downloader.Update(msg)

		// This one's the exception to the rule, moving to next step here because of progress bar's weird paradigm
		if m.downloader.Done {
			m.cleanServerFiles = NewCleanServerFilesModel(m.batchChoice)
			m.step = stepCleanServerFiles
			return m, m.cleanServerFiles.Init()
		}

		return m, cmd
	}

	return m, nil
}

func (m *Model) View() string {
	switch m.step {

	case stepEnvMenu:
		return m.peakEnvMenu.View()

	case stepIntegrationMenu:
		return m.integrationMenu.View()

	case stepBatchMenu:
		return m.batchMenu.View()

	case stepCopyingBatchFiles:
		return m.copyBatchFiles.View()

	case stepDownloading:
		return m.downloader.View()

	case stepCleanServerFiles:
		return m.cleanServerFiles.View()

	default:
		return "unknown state"
	}
}
