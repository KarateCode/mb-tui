package view_all_versions

import (
	"fmt"
	"image/color"
	"strings"

	tui "example.com/downloader/tui"
	downloader "example.com/downloader/tui/downloader"
	"github.com/charmbracelet/bubbles/spinner"
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

type getRequest struct {
	name       string
	version    string
	total      int64
	downloaded int64
	done       bool
}

type Model struct {
	getRequests []getRequest
	doneCount   int
	Done        bool

	step          step
	wwwincEnvMenu tui.MenuModel
	prefix        string

	integrationMenu  tui.MenuModel
	calcBatches      tui.SshCmdModel
	batchMenu        tui.MenuModel
	copyBatchFiles   tui.SshCmdModel
	downloader       downloader.Model
	cleanServerFiles tui.SshCmdModel

	batchChoice       string
	integrationChoice string
	spin              spinner.Model
	Program           *tea.Program
}

func (m *Model) Init() tea.Cmd {
	return m.spin.Tick
}

func NewViewAllVersionsModel() *Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205")) // magenta

	getRequests := make([]getRequest, len(links))
	return &Model{
		step:        0,
		getRequests: getRequests,
		spin:        s,
	}
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case ProgressGetVersion:
		d := &m.getRequests[msg.Index]
		d.version = msg.VersionString
		d.done = true
		m.doneCount += 1
		if m.doneCount >= len(links) {
			return m, tea.Quit
		}

		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.spin, cmd = m.spin.Update(msg)
	return m, cmd
}

func versionOrSpinner(m *Model, i int) string {
	if m.getRequests[i].done {
		return m.getRequests[i].version
	} else {
		return m.spin.View()
	}
}

func gradientText(text string, start, end color.RGBA) string {
	runes := []rune(text)
	n := len(runes)
	var out strings.Builder
	for i, r := range runes {
		// simple linear interpolation
		t := float64(i) / float64(n-1)
		rC := uint8(float64(start.R)*(1-t) + float64(end.R)*t)
		gC := uint8(float64(start.G)*(1-t) + float64(end.G)*t)
		bC := uint8(float64(start.B)*(1-t) + float64(end.B)*t)
		col := lipgloss.Color(fmt.Sprintf("#%02x%02x%02x", rC, gC, bC))
		out.WriteString(lipgloss.NewStyle().Foreground(col).Render(string(r)))
	}

	return out.String()
}

func (m *Model) View() string {
	width := 64 // or get from Bubble Tea window size messages
	title := lipgloss.Place(
		width,
		1,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.NewStyle().
			Bold(true).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")). // blue-ish
			Foreground(lipgloss.Color("#FFA500")).  // orange
			// Foreground(lipgloss.Color("#94E2D5")). // light teal
			Padding(0, 1).
			Render("  View All Versions  "),
	)

	i := 0
	subEnvStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("63")) // magenta
	header := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFA500")) // orange

	purple := color.RGBA{0x7C, 0x3A, 0xED, 0xFF}
	pink := color.RGBA{0xEC, 0x48, 0x99, 0xFF}

	var body string
	body += gradientText("\n\n=======================\t\t\t=======================\n", purple, pink)
	body += header.Render("Staging\t\t\t\t\tProduction")
	body += gradientText("\n=======================\t\t\t=======================\n", purple, pink)

	body += "\nWWWINC\t\t\t\t\tWWWINC\n"
	body += subEnvStyle.Render("NA\t\t")
	body += versionOrSpinner(m, i)
	i += 1
	body += subEnvStyle.Render("\t\t\tNA\t\t")
	body += versionOrSpinner(m, i)
	i += 1

	body += subEnvStyle.Render("\nEU\t\t")
	body += versionOrSpinner(m, i)
	i += 1
	body += subEnvStyle.Render("\t\t\tEU\t\t")
	body += versionOrSpinner(m, i)
	i += 1

	body += "\n\nPEAK\t\t\t\t\tPEAK\n"
	body += subEnvStyle.Render("NA\t\t")
	body += versionOrSpinner(m, i)
	i += 1
	body += subEnvStyle.Render("\t\t\tNA\t\t")
	body += versionOrSpinner(m, i)
	i += 1

	body += subEnvStyle.Render("\nEU\t\t")
	body += versionOrSpinner(m, i)
	i += 1
	body += subEnvStyle.Render("\t\t\tEU\t\t")
	body += versionOrSpinner(m, i)
	i += 1

	body += subEnvStyle.Render("\nCascade\t\t")
	body += versionOrSpinner(m, i)
	i += 1
	body += subEnvStyle.Render("\t\t\tCascade\t\t")
	body += versionOrSpinner(m, i)
	i += 1

	body += "\n\nConverse\t\t\t\tConverse\n"
	body += subEnvStyle.Render("NA\t\t")
	body += versionOrSpinner(m, i)
	i += 1
	body += subEnvStyle.Render("\t\t\tNA\t\t")
	body += versionOrSpinner(m, i)
	i += 1

	body += subEnvStyle.Render("\nEU\t\t")
	body += versionOrSpinner(m, i)
	i += 1
	body += subEnvStyle.Render("\t\t\tEU\t\t")
	body += versionOrSpinner(m, i)
	i += 1

	body += subEnvStyle.Render("\nLA\t\t")
	body += versionOrSpinner(m, i)
	i += 1
	body += subEnvStyle.Render("\t\t\tLA\t\t")
	body += versionOrSpinner(m, i)
	i += 1

	body += subEnvStyle.Render("\nAP\t\t")
	body += versionOrSpinner(m, i)
	i += 1
	body += subEnvStyle.Render("\t\t\tAP\t\t")
	body += versionOrSpinner(m, i)
	i += 1

	body += "\n\nCore and Enterprise Light\t\t\t\tCore and Enterprise Light"
	body += subEnvStyle.Render("\nCore\t\t")
	body += versionOrSpinner(m, i)
	i += 1
	body += subEnvStyle.Render("\t\t\tCore\t\t")
	body += versionOrSpinner(m, i)
	i += 1

	body += subEnvStyle.Render("\nNetsuite\t")
	body += versionOrSpinner(m, i)
	i += 1
	body += subEnvStyle.Render("\t\t\tNetsuite\t")
	body += versionOrSpinner(m, i)
	i += 1

	body += subEnvStyle.Render("\n\t\t\t\t\tDemo\t\t")
	body += versionOrSpinner(m, i)
	i += 1

	body += subEnvStyle.Render("\nCID\t\t")
	body += versionOrSpinner(m, i)
	i += 1
	body += subEnvStyle.Render("\t\t\tCID\t\t")
	body += versionOrSpinner(m, i)
	i += 1

	body += subEnvStyle.Render("\nDanpost\t\t")
	body += versionOrSpinner(m, i)
	i += 1
	body += subEnvStyle.Render("\t\t\tDanpost\t\t")
	body += versionOrSpinner(m, i)
	i += 1

	body += subEnvStyle.Render("\nOofos\t\t")
	body += versionOrSpinner(m, i)
	i += 1
	body += subEnvStyle.Render("\t\t\tOofos\t\t")
	body += versionOrSpinner(m, i)
	i += 1

	body += subEnvStyle.Render("\nVida\t\t")
	body += versionOrSpinner(m, i)
	i += 1
	body += subEnvStyle.Render("\t\t\tVida\t\t")
	body += versionOrSpinner(m, i)
	i += 1

	body += subEnvStyle.Render("\nLandau\t\t")
	body += versionOrSpinner(m, i)
	i += 1
	body += subEnvStyle.Render("\t\t\tLandau\t\t")
	body += versionOrSpinner(m, i)
	i += 1

	body += subEnvStyle.Render("\nJoules EU\t")
	body += versionOrSpinner(m, i)
	i += 1
	body += subEnvStyle.Render("\t\t\tJoules EU\t")
	body += versionOrSpinner(m, i)
	i += 1

	body += "\n\n"

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		body,
	)
}
