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
	message     string
	filterInput textinput.Model
	list        list.Model
	emitChoice  TeaCmdCallback

	selected string
}

var (
	tealStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("14")) // teal/cyan
	blueStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("12")) // blue
)

func NewMenu(message string, options []string, callback TeaCmdCallback) MenuModel {
	// message: 'Would you like to trim files to a certain import?',
	items := itemsFrom(options)

	// Text input
	ti := textinput.New()
	ti.Placeholder = "Search..."
	ti.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("11")) // yellow
	ti.Cursor.SetChar("|")
	ti.Focus()

	// List
	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = false
	delegate.SetHeight(1)
	delegate.SetSpacing(0)

	height := len(options) + 3
	l := list.New(items, delegate, 50, height) // WIDTH, HEIGHT
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.SetShowHelp(false)
	l.SetShowTitle(false)
	l.SetShowPagination(false)

	l.Styles.Title = lipgloss.NewStyle()

	return MenuModel{
		allOptions:  options,
		message:     message,
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

		case "ctrl+g":
			m.filterInput.SetValue("")
			return m, nil

		case "enter":
			if selected, ok := m.list.SelectedItem().(item); ok {
				m.selected = string(selected)
				teaCmd := m.emitChoice(m.selected)

				return m, teaCmd
			}
			return m, nil
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
	message := tealStyle.Render(m.message)
	// If you just want the prefix “> ” blue, but the typed text normal:
	// filter := m.filterInput.View()
	// if strings.HasPrefix(filter, ">") {
	// 	filter = blueStyle.Render(">") + filter[1:]
	// }

	// Color Reference:
	// "15" - White/text color
	// "14" — Cyan/Teal
	// "12" — Blue
	// "13" — Magenta
	// "11" — Yellow
	// "9" — Red
	filter := m.filterInput.View()
	if filter == ">" {
		filter = blueStyle.Render(">")
	} else {
		filter = blueStyle.Render(filter)
	}

	return fmt.Sprintf(
		"%s\n%s%s",
		message,
		filter,
		m.list.View(),
	)
}
