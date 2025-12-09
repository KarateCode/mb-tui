package peak_setup_integration

import (
	"bufio"
	"strings"

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
	stepDone
)

type integrationMenuComplete string
type batchMenuComplete string
type envMenuComplete string
type copyFilesComplete string
type cleanServerComplete string
type calcBatchesComplete string

type Model struct {
	step        step
	peakEnvMenu tui.MenuModel
	prefix      string

	integrationMenu  tui.MenuModel
	calcBatches      tui.SshCmdModel
	batchMenu        tui.MenuModel
	copyBatchFiles   tui.SshCmdModel
	downloader       downloader.Model
	cleanServerFiles tui.SshCmdModel

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
		"Select Peak Environment",
		names,
		func(selected string) tea.Cmd {
			teaCmd := func() tea.Msg {
				choice := envMenuComplete(selected)
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
	switch clientCode {
	case "cascade-na":
		return "lax_"
	case "bauer-na":
		return "hockey_na_"
	case "bauer-eu":
		return "hockey_eu_"
	default:
		return "hockey_eu_"
	}
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

	case envMenuComplete:
		selected := string(msg)

		foundEnv, _ := findEnvByName(environments(), selected)
		m.envMenuChoice = foundEnv
		m.prefix = calcPrefix(foundEnv.clientCode)
		m.integrationMenu = tui.NewMenu(
			"Which Integration batch would you like?",
			integrationListItems(),
			func(selected string) tea.Cmd {
				teaCmd := func() tea.Msg {
					return integrationMenuComplete(selected)
				}
				return teaCmd
			},
		)
		m.step = stepIntegrationMenu
		return m, m.integrationMenu.Init()

	case integrationMenuComplete:
		m.step = stepCalcBatches
		choice := string(msg)
		showBatchesCmd := commandForIntegration(choice, m.envMenuChoice)

		m.calcBatches = tui.NewShellCmd(
			m.envMenuChoice.sshServer,
			showBatchesCmd,
			"Calculating batches...",
			func(output string) tea.Msg {
				return calcBatchesComplete(output)
			},
		)
		return m, m.calcBatches.Init()

	case calcBatchesComplete:
		output := calcBatchesComplete(msg)
		lines := linesFromOutput(string(output))

		m.batchMenu = tui.NewMenu(
			"Please select batch number",
			lines,
			func(selected string) tea.Cmd {
				teaCmd := func() tea.Msg {
					choice := batchMenuComplete(selected)
					return choice
				}
				return teaCmd
			},
		)
		m.step = stepBatchMenu
		return m, m.batchMenu.Init()

	case batchMenuComplete:
		choice := string(msg)
		m.batchChoice = choice
		copyFilesCmd := generateCopyFilesCmd(m.envMenuChoice, choice)

		m.copyBatchFiles = tui.NewShellCmd(
			m.envMenuChoice.sshServer,
			copyFilesCmd,
			"Copying files to /client/dumps...",
			func(output string) tea.Msg {
				return copyFilesComplete(output)
			},
		)
		m.step = stepCopyingBatchFiles
		return m, m.copyBatchFiles.Init()

	case copyFilesComplete:
		output := copyFilesComplete(msg)
		fileNames := linesFromOutput(string(output))

		m.downloader = downloader.NewModel(fileNames)
		downloadFiles := func() tea.Msg {
			return downloader.DownloadFiles(fileNames, m.envMenuChoice.sshServer, m.Program)
		}
		m.step = stepDownloading
		return m, downloadFiles

	case cleanServerComplete:
		m.step = stepDone
		return m, nil

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

		// This one that's the exception to the rule; moving to the next step here because of progress bar's weird paradigm
		if m.downloader.Done {
			copyFilesCmd := generateCleanServerFilesCmd(m.envMenuChoice, m.batchChoice)
			m.cleanServerFiles = tui.NewShellCmd(
				m.envMenuChoice.sshServer,
				copyFilesCmd,
				"Cleaning up files on the server...",
				func(output string) tea.Msg {
					return cleanServerComplete(output)
				},
			)
			m.step = stepCleanServerFiles
			return m, m.cleanServerFiles.Init()
		}
		return m, cmd

	case stepCleanServerFiles:
		m.cleanServerFiles, cmd = m.cleanServerFiles.Update(msg)
		return m, cmd

	case stepDone:
		return m, tea.Quit
	}

	return m, nil
}

func completedTodosView(m *Model) string {
	greenStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#8aff8a")).Italic(true)
	completedTodos := ""
	if m.step > stepEnvMenu {
		completedTodos += "\n  ✔ Peak environment chosen"
	}
	if m.step > stepIntegrationMenu {
		completedTodos += "\n  ✔ Integration type selected"
	}
	if m.step > stepCalcBatches {
		completedTodos += "\n  ✔ Batches calculated"
	}
	if m.step > stepBatchMenu {
		completedTodos += "\n  ✔ Batch " + m.batchChoice + " chosen"
	}
	if m.step > stepCopyingBatchFiles {
		completedTodos += "\n  ✔ Files copied to /client/dumps"
	}
	if m.step > stepDownloading {
		completedTodos += "\n  ✔ All files in batch " + m.batchChoice + " downloaded"
	}
	if m.step > stepCleanServerFiles {
		completedTodos += "\n  ✔ Cleaned out files in /client/dumps"
	}
	completedTodos += "\n\n"

	return greenStyle.Render(completedTodos)
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

	completedTodos := completedTodosView(m)

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

	case stepDone:
		body = ""

	default:
		body = "unknown state"
	}

	// 3. Stack title + body vertically
	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		completedTodos,
		body,
	)
}
