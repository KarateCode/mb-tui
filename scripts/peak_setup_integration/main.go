package peak_setup_integration

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

func PeakSetupIntegration() {
	model := NewModel()
	program := tea.NewProgram(model)
	model.Program = program
	if _, err := program.Run(); err == nil {
		fmt.Println("exiting p.Run")
	}
}
