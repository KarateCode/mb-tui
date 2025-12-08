package peak_setup_integration

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Can we resue all this as just a 'menuItem' type for all my menus?
// type batchItem string
// func (i batchItem) Title() string       { return string(i) }
// func (i batchItem) Description() string { return "" }
// func (i batchItem) FilterValue() string { return string(i) }

// type (
// 	downloadCompleteMsg []string
// )

type EnvMenuModel struct {
	allEnvs     peakEnvs
	filterInput textinput.Model
	list        list.Model
	allOptions  []string

	isDownloading    bool
	downloadComplete bool
	showBatchesCmd   string

	Done     bool
	quitting bool
	selected string
}

type peakEnv struct {
	name       string
	sshServer  string
	subFolder  string
	clientCode string
}
type peakEnvs []peakEnv

func NewEnvMenu() EnvMenuModel {
	environments := peakEnvs{
		{
			name:       "Bauer EU Staging",
			sshServer:  "bauer-stag-eu-cf-integration",
			subFolder:  "EU",
			clientCode: "bauer-eu",
		},
		{
			name:       "Bauer EU Production",
			sshServer:  "bauer-prod-eu-cf-integration",
			subFolder:  "EU",
			clientCode: "bauer-eu",
		},
		{
			name:       "Bauer NA Staging",
			sshServer:  "bauer-stag-na-cf-integration",
			subFolder:  "NA",
			clientCode: "bauer-na",
		},
		{
			name:       "Bauer NA Production",
			sshServer:  "bauer-prod-na-cf-integration",
			subFolder:  "NA",
			clientCode: "bauer-na",
		},
		{
			name:       "Cascade NA Staging",
			sshServer:  "cascade-stag-na-cf-integration",
			subFolder:  "NA",
			clientCode: "cascade-na",
		},
		{
			name:       "Cascade NA Production",
			sshServer:  "cascade-prod-na-cf-integration",
			subFolder:  "NA",
			clientCode: "cascade-na",
		},
	}

	// Convert to list items
	items := make([]list.Item, len(environments))
	names := make([]string, len(environments))
	for i, s := range environments {
		items[i] = item(s.name)
		names[i] = s.name
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

	l := list.New(items, delegate, 50, 10) // WIDTH, HEIGHT
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.SetShowHelp(false)
	l.SetShowTitle(false)
	l.SetShowPagination(false)

	l.Styles.Title = lipgloss.NewStyle()

	return EnvMenuModel{
		allEnvs:     environments,
		allOptions:  names,
		filterInput: ti,
		list:        l,
		Done:        false,
	}
}

func (m EnvMenuModel) Init() tea.Cmd {
	return nil
	// return doDownload(m)
}

func (m EnvMenuModel) Update(msg tea.Msg) (EnvMenuModel, tea.Cmd) {
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
			selected := fmt.Sprintf("%+v", m.list.SelectedItem())
			m.selected = string(selected)

			foundEnv, found := findEnvByName(m.allEnvs, selected)

			if found {
				// fmt.Printf("Found environment: %+v\n", foundEnv)
				teaCmd := func() tea.Msg {
					choice := EnvMenuChoice(foundEnv)
					return choice
				}
				return m, teaCmd
				// } else {
				// 	fmt.Printf("Environment with name '%s' not found.\n", selected)
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
				filtered = append(filtered, batchItem(b))
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

// Define a function that mimics _.find behavior
func findEnvByName(envs peakEnvs, selectedName string) (peakEnv, bool) {
	for _, env := range envs {
		if env.name == selectedName {
			return env, true
		}
	}

	return peakEnv{}, false
}

func (m EnvMenuModel) View() string {
	if m.quitting {
		return ""
	}

	return fmt.Sprintf(
		"Please select which Peak environment to download from\nFilter: %s\n\n%s",
		m.filterInput.View(),
		m.list.View(),
	)
}

func (m EnvMenuModel) Selected() string {
	return m.selected
}
