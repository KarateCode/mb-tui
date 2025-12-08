package peak_setup_integration

import (
	// "fmt"
	// "strings"

	// "example.com/downloader/batchmenu"
	// "github.com/charmbracelet/bubbles/progress"
	// batchmenu "example.com/downloader/batchmenu"
	// "fmt"

	downloader "example.com/downloader/tui/downloader"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	prefix           string
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

func calcPrefix(clientCode string) string {
	if clientCode == "cascade-na" {
		return "lax_"
	} else if clientCode == "bauer-na" {
		return "hockey_na_"
	} else if clientCode == "bauer-eu" {
		return "hockey_eu_"
	}
	return "hockey_eu_"
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case EnvMenuChoice:
		choice := peakEnv(msg)
		m.envMenuChoice = choice
		m.prefix = calcPrefix(choice.clientCode)
		m.integrationMenu = NewIntegrationMenu()
		m.step = stepIntegrationMenu
		return m, m.integrationMenu.Init()

	case IntegrationMenuChoice:
		choice := IntegrationMenuChoice(msg)
		m.batchMenu = NewMenu(choice, m.envMenuChoice)
		m.step = stepBatchMenu
		return m, m.batchMenu.Init()

	case BatchChoice:
		choice := string(msg)
		m.batchChoice = choice
		m.copyBatchFiles = NewCopyBatchFilesModel(choice, m.envMenuChoice)
		m.step = stepCopyingBatchFiles
		return m, m.copyBatchFiles.Init()

	case copyCompleteMsg:
		fileNames := []string(msg)
		m.downloader = downloader.NewModel(fileNames)
		downloadFiles := func() tea.Msg {
			return downloader.DownloadFiles(fileNames, m.envMenuChoice.sshServer, m.Program)
		}
		m.step = stepDownloading
		return m, downloadFiles

	case cleanServerCompleteMsg:
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
			m.cleanServerFiles = NewCleanServerFilesModel(m.batchChoice, m.envMenuChoice)
			m.step = stepCleanServerFiles
			return m, m.cleanServerFiles.Init()
		}

		return m, cmd
	}

	return m, nil
}

func (m *Model) View() string {
	width := 80 // or get from Bubble Tea window size messages
	title := lipgloss.Place(
		width,
		1,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.NewStyle().
			Bold(true).
			Border(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.Color("63")). // blue-ish
			Foreground(lipgloss.Color("#FFA500")).  // orange
			// Foreground(lipgloss.Color("#94E2D5")). // light teal
			Padding(0, 1).
			Render("  Peak Integration Setup  "),
	)

	// 2. Render active submodel
	var body string
	switch m.step {

	case stepEnvMenu:
		body = m.peakEnvMenu.View()

	case stepIntegrationMenu:
		body = m.integrationMenu.View()

	case stepBatchMenu:
		body = m.batchMenu.View()

	case stepCopyingBatchFiles:
		body = m.copyBatchFiles.View()

	case stepDownloading:
		body = m.downloader.View()

	case stepCleanServerFiles:
		body = m.cleanServerFiles.View()

	default:
		body = "unknown state"
	}

	// 3. Stack title + body vertically
	return lipgloss.JoinVertical(lipgloss.Left, title, body)
}
