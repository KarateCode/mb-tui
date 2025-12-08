package tui

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

type TeaCmdCallback func(selected string) tea.Cmd

type MenuModel struct {
	allOptions  []string
	filterInput textinput.Model
	list        list.Model
	emitChoice  TeaCmdCallback

	selected string
}

func NewMenu(options []string, callback TeaCmdCallback) MenuModel {
	// message: 'Would you like to trim files to a certain import?',
	items := itemsFrom(options)

	// Text input
	ti := textinput.New()
	ti.Placeholder = "Search..."
	ti.Focus()

	// List
	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = false
	delegate.SetHeight(1)
	delegate.SetSpacing(0)

	l := list.New(items, delegate, 50, 40) // WIDTH, HEIGHT
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.SetShowHelp(false)
	l.SetShowTitle(false)
	l.SetShowPagination(false)

	l.Styles.Title = lipgloss.NewStyle()

	return MenuModel{
		allOptions:  options,
		filterInput: ti,
		list:        l,
		emitChoice:  callback,
	}

}

func (m MenuModel) Init() tea.Cmd {
	return nil
}

func itemsFrom(options []string) []list.Item {
	items := make([]list.Item, 0, len(options))
	for _, b := range options {
		items = append(items, item(b))
	}
	return items
}

func (m MenuModel) Update(msg tea.Msg) (MenuModel, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "ctrl+n":
			m.list.CursorDown()
			return m, nil

		case "ctrl+p":
			m.list.CursorUp()
			return m, nil

		case "enter":
			if selected, ok := m.list.SelectedItem().(item); ok {
				m.selected = string(selected)
				teaCmd := m.emitChoice(m.selected)
				// teaCmd := func() tea.Msg {
				// 	choice := IntegrationMenuChoice(m.selected)
				// 	return choice
				// }

				return m, teaCmd
			}
			return m, nil // should probably send back a Cmd where tea.Msg is the selected value
			// might have to convert from item to string
			// IntegrationMenuChoice
		}
	}

	// Update text input
	var cmd tea.Cmd
	m.filterInput, cmd = m.filterInput.Update(msg)

	// Filter list items
	filter := strings.ToLower(m.filterInput.Value())

	if filter != "" {
		var filtered []list.Item
		for _, b := range m.allOptions {
			if strings.Contains(strings.ToLower(b), filter) {
				filtered = append(filtered, item(b))
			}
		}
		m.list.SetItems(filtered)
	} else {
		// Only reset when returning to full list
		m.list.SetItems(itemsFrom(m.allOptions))
	}

	// Update list
	m.list, _ = m.list.Update(msg)

	return m, cmd
}

func (m MenuModel) View() string {
	return fmt.Sprintf(
		"Filter: %s\n\n%s",
		m.filterInput.View(),
		m.list.View(),
	)
}
