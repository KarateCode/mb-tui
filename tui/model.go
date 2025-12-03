package tui

import (
	// "fmt"
	"strings"

	// "example.com/downloader/batchmenu"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

type fileDownload struct {
	total      int64
	downloaded int64
	bar        progress.Model
	done       bool
}
type model struct {
	fileDownloads []fileDownload
	doneCount     int
}

type setTotalMsg struct {
	Index int
	Total int64
}

type progressMsg struct {
	Index int
	Bytes int64
}

type doneMsg struct {
	Index int
}

func newModel(fileNames []string) model {
	fileCount := len(fileNames)
	m := model{
		fileDownloads: make([]fileDownload, fileCount),
	}

	for i := range m.fileDownloads {
		m.fileDownloads[i].bar = progress.New(progress.WithDefaultGradient())
	}

	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// fmt.Printf("from update: %+v\n", msg)
	switch msg := msg.(type) {

	case progressMsg:
		d := &m.fileDownloads[msg.Index]
		d.downloaded = msg.Bytes

		if d.downloaded >= d.total {
			d.done = true
			// return m, tea.Quit
		}
		return m, nil

	case setTotalMsg:
		// fmt.Printf("\n\n\n\n\nsetTotalMsg: %+v %+v \n", msg.Index, msg.Total)
		m.fileDownloads[msg.Index].total = msg.Total
		return m, nil

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

	case doneMsg:
		m.doneCount++
		if m.doneCount == len(m.fileDownloads) {
			return m, tea.Quit
		}
		return m, nil
	}
	return m, nil
}

func (m model) View() string {
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
