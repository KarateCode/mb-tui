package view_all_versions

import (
	tea "github.com/charmbracelet/bubbletea"
)

func ViewAllVersions() {
	model := NewViewAllVersionsModel()
	program := tea.NewProgram(model)
	model.Program = program
	if _, err := program.Run(); err == nil {
		// fmt.Println("exiting p.Run")
	}
}
