package main

import (
	viewAllVersions "example.com/downloader/scripts/view_all_versions"
	wwwSetupIntegration "example.com/downloader/scripts/wwwinc_setup_integration"
	tui "example.com/downloader/tui"
	tea "github.com/charmbracelet/bubbletea"
)

type step int
type scriptChoiceComplete string

const (
	stepScriptMenu step = iota
	stepViewAllVersions
	stepWwwSetupIntegration
)

type Model struct {
	step step

	scriptMenu               tui.MenuModel
	viewAllVersionsModel     *viewAllVersions.Model
	wwwSetupIntegrationModel *wwwSetupIntegration.Model
	Program                  *tea.Program
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func NewTopModel() *Model {
	menuItems := []string{
		"View All Versions",
		"Wwwinc Setup Integration",
		"Peak Setup Integration",
	}

	m := tui.NewMenu(
		"Select Script",
		menuItems,
		func(selected string) tea.Cmd {
			teaCmd := func() tea.Msg {
				choice := scriptChoiceComplete(selected)
				return choice
			}
			return teaCmd
		},
	)

	return &Model{
		step:       0,
		scriptMenu: m,
	}
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case scriptChoiceComplete:
		choice := scriptChoiceComplete(msg)
		if choice == "View All Versions" {
			m.step = stepViewAllVersions
			m.viewAllVersionsModel = viewAllVersions.NewViewAllVersionsModel()
			teaCmd := func() tea.Msg {
				return viewAllVersions.GetVersionsOverHttp(m.Program)
			}
			return m, teaCmd
		} else if choice == "Wwwinc Setup Integration" {
			m.step = stepWwwSetupIntegration
			model := wwwSetupIntegration.NewModel()
			model.Program = m.Program
			m.wwwSetupIntegrationModel = model
			teaCmd := func() tea.Msg {
				return viewAllVersions.GetVersionsOverHttp(m.Program)
			}
			return m, teaCmd
		}

		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		}
	}

	switch m.step {
	case stepScriptMenu:
		var cmd tea.Cmd
		m.scriptMenu, cmd = m.scriptMenu.Update(msg)
		return m, cmd

	case stepViewAllVersions:
		var cmd tea.Cmd
		_, cmd = m.viewAllVersionsModel.Update(msg)
		// m.viewAllVersionsModel = &model
		return m, cmd

	case stepWwwSetupIntegration:
		var cmd tea.Cmd
		_, cmd = m.wwwSetupIntegrationModel.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m *Model) View() string {
	var body string

	switch m.step {
	case stepScriptMenu:
		body = m.scriptMenu.View()
	case stepViewAllVersions:
		body = m.viewAllVersionsModel.View()
	case stepWwwSetupIntegration:
		body = m.wwwSetupIntegrationModel.View()

	default:
		body = "unknown state"
	}

	return body
}
