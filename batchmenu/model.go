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
	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = false
	delegate.SetHeight(1)
	delegate.SetSpacing(0)

	l := list.New(items, delegate, 50, 40) // WIDTH=50, HEIGHT=20 rows
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.SetShowHelp(false)
	l.SetShowTitle(false)
	l.SetShowPagination(false)

	l.Styles.Title = lipgloss.NewStyle()

	return model{
		allBatches:  batchList,
		filterInput: ti,
		list:        l,
	}
}

func (m model) Init() tea.Cmd {
	return nil
	// return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit

			// Emacs-style movement
		case "ctrl+n":
			// Move down
			m.list.CursorDown()
			return m, nil

		case "ctrl+p":
			// Move up
			m.list.CursorUp()
			return m, nil

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

	// var filtered []list.Item
	// for _, b := range m.allBatches {
	// 	if strings.Contains(strings.ToLower(b), filter) {
	// 		filtered = append(filtered, item(b))
	// 	}
	// }
	// m.list.SetItems(filtered)
	if filter != "" {
		var filtered []list.Item
		for _, b := range m.allBatches {
			if strings.Contains(strings.ToLower(b), filter) {
				filtered = append(filtered, item(b))
			}
		}
		m.list.SetItems(filtered)
	} else {
		// Only reset when returning to full list
		m.list.SetItems(itemsFrom(m.allBatches))
	}

	// Update list
	m.list, _ = m.list.Update(msg)

	return m, cmd
}

// itemsFrom converts []string -> []list.Item
func itemsFrom(batches []string) []list.Item {
	items := make([]list.Item, 0, len(batches))
	for _, b := range batches {
		items = append(items, item(b))
	}
	return items
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
