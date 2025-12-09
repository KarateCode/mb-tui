package downloader

import (
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

type fileDownload struct {
	total      int64
	downloaded int64
	bar        progress.Model
	done       bool
}
type Model struct {
	fileDownloads []fileDownload
	doneCount     int
	Done          bool
}

type SetTotalMsg struct {
	Index int
	Total int64
}

type ProgressMsg struct {
	Index int
	Bytes int64
}

type DoneMsg struct {
	Index int
}

func NewModel(fileNames []string) Model {
	fileCount := len(fileNames)
	m := Model{
		fileDownloads: make([]fileDownload, fileCount),
	}

	for i := range m.fileDownloads {
		m.fileDownloads[i].bar = progress.New(progress.WithDefaultGradient())
	}

	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	// fmt.Printf("from update: %+v\n", msg)
	switch msg := msg.(type) {

	case ProgressMsg:
		d := &m.fileDownloads[msg.Index]
		d.downloaded = msg.Bytes

		if d.downloaded >= d.total {
			d.done = true
		}
		return m, nil

	case SetTotalMsg:
		m.fileDownloads[msg.Index].total = msg.Total
		return m, nil

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

	case DoneMsg:
		m.doneCount++
		if m.doneCount == len(m.fileDownloads) {
			m.Done = true
			return m, tea.Quit
		}
		return m, nil
	}
	return m, nil
}

func (m Model) View() string {
	var lines []string
	lines = append(lines, "\n", "")

	for _, d := range m.fileDownloads {
		var ratio float64
		if d.total > 0 {
			ratio = float64(d.downloaded) / float64(d.total)
		}

		lines = append(lines, d.bar.ViewAs(ratio), "")
	}

	return strings.Join(lines, "\n")
}
