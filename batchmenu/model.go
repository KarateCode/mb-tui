package batchmenu

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type item string

func (i item) Title() string       { return string(i) }
func (i item) Description() string { return "" }
func (i item) FilterValue() string { return string(i) }

type model struct {
	allBatches  []string
	filterInput textinput.Model
	list        list.Model

	quitting bool
	selected string
}

func NewMenu(batchList []string) model {
	// Convert to list items
	items := make([]list.Item, len(batchList))
	for i, s := range batchList {
		items[i] = item(s)
	}

	// Text input
	ti := textinput.New()
	ti.Placeholder = "Search..."
	ti.Focus()

	// List
	l := list.New(items, list.NewDefaultDelegate(), 50, 20) // WIDTH=50, HEIGHT=20 rows
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.SetShowHelp(false)
	l.SetShowTitle(false)
	l.SetShowPagination(false)

	// Remove all padding/margins for tight rows
	// l.Styles.NormalTitle = lipgloss.NewStyle()
	l.Styles.Title = lipgloss.NewStyle()
	l.SetDelegate(newCompactDelegate())

	return model{
		allBatches:  batchList,
		filterInput: ti,
		list:        l,
	}
}

func newCompactDelegate() list.DefaultDelegate {
	d := list.NewDefaultDelegate()
	// d.Styles.Normal = lipgloss.NewStyle()
	// d.Styles.Selected = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	return d
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			if selected, ok := m.list.SelectedItem().(item); ok {
				m.selected = string(selected)
			}
			return m, tea.Quit
		}
	}

	// Update text input
	var cmd tea.Cmd
	m.filterInput, cmd = m.filterInput.Update(msg)

	// Filter list items
	filter := strings.ToLower(m.filterInput.Value())

	var filtered []list.Item
	for _, b := range m.allBatches {
		if strings.Contains(strings.ToLower(b), filter) {
			filtered = append(filtered, item(b))
		}
	}

	m.list.SetItems(filtered)

	// Update list
	m.list, _ = m.list.Update(msg)

	return m, cmd
}

func (m model) View() string {
	if m.quitting {
		return ""
	}

	return fmt.Sprintf(
		"Filter: %s\n\n%s",
		m.filterInput.View(),
		m.list.View(),
	)
}

func (m model) Selected() string {
	return m.selected
}
