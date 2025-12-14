package wwwinc_setup_integration

import (
	tea "github.com/charmbracelet/bubbletea"
)

func WwwincSetupIntegration() {
	model := newModel()
	program := tea.NewProgram(model)
	model.Program = program
	if _, err := program.Run(); err == nil {
		// fmt.Println("exiting p.Run")
	}
}
