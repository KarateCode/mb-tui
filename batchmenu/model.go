package batchmenu

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
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
	items := make([]list.Item, len(batchList))
	for i, b := range batchList {
		items[i] = item(b)
	}

	ti := textinput.New()
	ti.Placeholder = "filter..."
	ti.Focus()

	ls := list.New(items, list.NewDefaultDelegate(), 30, 15)
	ls.Title = "Which Batch?"

	return model{
		allBatches:  batchList,
		filterInput: ti,
		list:        ls,
	}
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
