package main

import (
	viewAllVersions "example.com/downloader/scripts/view_all_versions"
	tui "example.com/downloader/tui"
	tea "github.com/charmbracelet/bubbletea"
)

type step int
type scriptChoiceComplete string

const (
	stepScriptMenu step = iota
	stepExecute
)

type Model struct {
	step step

	scriptMenu           tui.MenuModel
	viewAllVersionsModel *viewAllVersions.Model
	Program              *tea.Program
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
			m.step = stepExecute
			m.viewAllVersionsModel = viewAllVersions.NewViewAllVersionsModel()
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

	case stepExecute:
		var cmd tea.Cmd
		_, cmd = m.viewAllVersionsModel.Update(msg)
		// m.viewAllVersionsModel = &model
		return m, cmd
	}

	return m, nil
}

func (m *Model) View() string {
	var body string

	switch m.step {
	case stepScriptMenu:
		body = m.scriptMenu.View()
	case stepExecute:
		body = m.viewAllVersionsModel.View()

	default:
		body = "unknown state"
	}

	return body
}
