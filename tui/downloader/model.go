package downloader

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type fileDownload struct {
	name       string
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

	for i, fileName := range fileNames {
		m.fileDownloads[i].name = fileName
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
	grayItalic := lipgloss.NewStyle().Foreground(lipgloss.Color("245")).Italic(true)

	for _, download := range m.fileDownloads {
		var ratio float64
		if download.total > 0 {
			ratio = float64(download.downloaded) / float64(download.total)
		}

		line := fmt.Sprintf(
			"%s\n%s",
			grayItalic.Render(download.name),
			download.bar.ViewAs(ratio),
		)
		lines = append(lines, line, "")
	}

	return strings.Join(lines, "\n")
}
