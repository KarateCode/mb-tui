package batchmenu

import (
	"fmt"
	"io"
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

type singleLineItemDelegate struct{}

func (d singleLineItemDelegate) Height() int                               { return 1 }
func (d singleLineItemDelegate) Spacing() int                              { return 0 }
func (d singleLineItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

func (d singleLineItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := string(i)

	if index == m.Index() {
		str = "> " + str
	} else {
		str = "  " + str
	}

	fmt.Fprint(w, str)
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

	// delegate := singleLineItemDelegate{}
	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = false
	delegate.SetHeight(1)
	delegate.SetSpacing(0)
	// delegate.Styles.NormalTitle = delegate.Styles.NormalTitle.Padding(10, 10, 10, 11)
	// delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.Padding(10, 10, 10, 11)

	l := list.New(items, delegate, 50, 40) // WIDTH=50, HEIGHT=20 rows
	// l.SetShowDescription(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.SetShowHelp(false)
	l.SetShowTitle(false)
	l.SetShowPagination(false)

	// Remove all padding/margins for tight rows
	// l.Styles.NormalTitle = lipgloss.NewStyle()
	l.Styles.Title = lipgloss.NewStyle()
	// l.SetDelegate(newCompactDelegate())

	return model{
		allBatches:  batchList,
		filterInput: ti,
		list:        l,
	}
}

func newCompactDelegate() list.DefaultDelegate {
	d := list.NewDefaultDelegate()
	d.ShowDescription = false
	d.SetHeight(4)
	// d.Styles.NormalTitle = d.Styles.NormalTitle.Padding(0, 0, 0, 3)
	d.Styles.NormalTitle = d.Styles.NormalTitle.Padding(0, 0, 0, 1)
	d.Styles.SelectedTitle = d.Styles.SelectedTitle.Padding(0, 0, 0, 1)
	// d.Styles.NormalText = d.Styles.NormalText.Padding(0)
	// d.Styles.SelectedText = d.Styles.SelectedText.Padding(0)
	// d.Height = 1
	// d.Styles.SelectedTitle = d.Styles.SelectedTitle.Padding(0, 0, 0, 1)
	// d.Styles.NormalTitle = lipgloss.NewStyle().Foreground(lipgloss.Color("99")).MarginBottom(0).MarginTop(0)
	// d.Styles.SelectedTitle = lipgloss.NewStyle().Foreground(lipgloss.Color("212")).MarginBottom(0)
	// marginBottom := d.Styles
	// fmt.Printf("marginBottom:\n")
	// fmt.Printf("%+v\n", marginBottom)
	return d
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
