package peak_setup_integration

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type batchItem string

func (i batchItem) Title() string       { return string(i) }
func (i batchItem) Description() string { return "" }
func (i batchItem) FilterValue() string { return string(i) }

type (
	calcBatchesCompleteMsg []string
)

type BatchModel struct {
	allBatches  []string
	filterInput textinput.Model
	list        list.Model

	isDownloading          bool
	calcBatchesCompleteMsg bool
	showBatchesCmd         string

	Done     bool
	quitting bool
	selected string
}

func NewMenu(lines calcBatchesCompleteMsg) BatchModel {
	// Convert to list items
	items := make([]list.Item, len(lines))
	for i, s := range lines {
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

	return BatchModel{
		allBatches:  lines,
		filterInput: ti,
		list:        l,
		Done:        false,
	}
}

func (m BatchModel) Init() tea.Cmd {
	return nil
}

func (m BatchModel) Update(msg tea.Msg) (BatchModel, tea.Cmd) {
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
			selected := fmt.Sprintf("%+v", m.list.SelectedItem())
			m.selected = string(selected)

			teaCmd := func() tea.Msg {
				choice := BatchChoice(m.selected)
				return choice
			}

			return m, teaCmd
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
				filtered = append(filtered, batchItem(b))
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

func (m BatchModel) View() string {
	if m.quitting {
		return ""
	}

	return fmt.Sprintf(
		"\n\n✅ Download complete! Press Ctrl+C to exit.\nFilter: %s\n\n%s",
		m.filterInput.View(),
		m.list.View(),
	)
}

func (m BatchModel) Selected() string {
	return m.selected
}
