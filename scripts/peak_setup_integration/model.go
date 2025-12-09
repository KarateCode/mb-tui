package peak_setup_integration

import (
	"bufio"
	"fmt"
	"strings"
	// "example.com/downloader/batchmenu"
	// "github.com/charmbracelet/bubbles/progress"
	// batchmenu "example.com/downloader/batchmenu"
	// "fmt"

	tui "example.com/downloader/tui"
	downloader "example.com/downloader/tui/downloader"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type step int

const (
	stepEnvMenu step = iota
	stepIntegrationMenu
	stepCalcBatches
	stepBatchMenu
	stepCopyingBatchFiles
	stepDownloading
	stepCleanServerFiles
)

type IntegrationMenuChoice string
type BatchChoice string
type EnvMenuChoice string
type copyCompleteMsg []string
type cleanServerCompleteMsg string
type calcBatchesCompleteMsg string

type Model struct {
	step        step
	peakEnvMenu tui.MenuModel
	prefix      string

	integrationMenu  tui.MenuModel
	calcBatches      tui.SshCmdModel
	batchMenu        tui.MenuModel
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
	envs := environments()
	names := make([]string, len(envs))
	for i, s := range envs {
		names[i] = s.name
	}

	m := tui.NewMenu(
		names,
		func(selected string) tea.Cmd {
			teaCmd := func() tea.Msg {
				choice := EnvMenuChoice(selected)
				return choice
			}
			return teaCmd
		},
	)

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

func linesFromOutput(output string) []string {
	lines := []string{}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}

	return lines
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case EnvMenuChoice:
		selected := string(msg)

		foundEnv, _ := findEnvByName(environments(), selected)
		m.envMenuChoice = foundEnv
		m.prefix = calcPrefix(foundEnv.clientCode)
		m.integrationMenu = tui.NewMenu(
			[]string{
				"Nope! Give me them all",
				"Product Import",
				"Customer Import",
				"Inventory Import",
				"SalesRep Import",
				"BG/BHC import",
				"SalesOrg/PoType Import",
			},
			func(selected string) tea.Cmd {
				teaCmd := func() tea.Msg {
					return IntegrationMenuChoice(selected)
				}
				return teaCmd
			},
		)
		m.step = stepIntegrationMenu
		return m, m.integrationMenu.Init()

	case IntegrationMenuChoice:
		m.step = stepCalcBatches
		choice := string(msg)
		showBatchesCmd := commandForIntegration(choice, m.envMenuChoice)

		m.calcBatches = tui.NewShellCmd(
			m.envMenuChoice.sshServer,
			showBatchesCmd,
			"Calculating batches from stuff...",
			func(output string) tea.Msg {
				return calcBatchesCompleteMsg(output)
			},
		)
		return m, m.calcBatches.Init()

	case calcBatchesCompleteMsg:
		output := calcBatchesCompleteMsg(msg)
		lines := linesFromOutput(string(output))

		m.batchMenu = tui.NewMenu(
			lines,
			func(selected string) tea.Cmd {
				teaCmd := func() tea.Msg {
					choice := BatchChoice(selected)
					return choice
				}
				return teaCmd
			},
		)
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
		fmt.Println("Clean exit")
		return m, tea.Quit

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	switch m.step {

	case stepEnvMenu:
		m.peakEnvMenu, cmd = m.peakEnvMenu.Update(msg)
		return m, cmd

	case stepIntegrationMenu:
		m.integrationMenu, cmd = m.integrationMenu.Update(msg)
		return m, cmd

	case stepCalcBatches:
		m.calcBatches, cmd = m.calcBatches.Update(msg)
		return m, cmd

	case stepBatchMenu:
		m.batchMenu, cmd = m.batchMenu.Update(msg)
		return m, cmd

	case stepCopyingBatchFiles:
		m.copyBatchFiles, cmd = m.copyBatchFiles.Update(msg)
		return m, cmd

	case stepDownloading:
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

	case stepCalcBatches:
		body = m.calcBatches.View()

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
