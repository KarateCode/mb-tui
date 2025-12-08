package peak_setup_integration

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

type IntegrationMenuModel struct {
	allBatches  []string
	filterInput textinput.Model
	list        list.Model

	Done     bool
	quitting bool
	selected string
}

func NewIntegrationMenu() IntegrationMenuModel {
	// message: 'Would you like to trim files to a certain import?',
	fileTypes := []string{
		"Nope! Give me them all",
		"Product Import",
		"Customer Import",
		"Inventory Import",
		"SalesRep Import",
		"BG/BHC import",
		"SalesOrg/PoType Import",
	}
	items := make([]list.Item, len(fileTypes))
	for i, s := range fileTypes {
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

	l := list.New(items, delegate, 50, 40) // WIDTH, HEIGHT
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.SetShowHelp(false)
	l.SetShowTitle(false)
	l.SetShowPagination(false)

	l.Styles.Title = lipgloss.NewStyle()

	return IntegrationMenuModel{
		allBatches:  fileTypes,
		filterInput: ti,
		list:        l,
		Done:        false,
	}

}

func (m IntegrationMenuModel) Init() tea.Cmd {
	return nil
}

// itemsFrom converts []string -> []list.Item
func itemsFrom(batches []string) []list.Item {
	items := make([]list.Item, 0, len(batches))
	for _, b := range batches {
		items = append(items, item(b))
	}
	return items
}

func (m IntegrationMenuModel) Update(msg tea.Msg) (IntegrationMenuModel, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.quitting = true
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
				// m.Done = true
				teaCmd := func() tea.Msg {
					choice := IntegrationMenuChoice(m.selected)
					return choice
				}

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

func (m IntegrationMenuModel) View() string {
	if m.quitting {
		return ""
	}

	return fmt.Sprintf(
		"Filter: %s\n\n%s",
		m.filterInput.View(),
		m.list.View(),
	)
}
