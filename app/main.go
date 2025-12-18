package main

import (
	// peakSetupIntegration "example.com/downloader/scripts/peak_setup_integration"
	// wwwincSetupIntegration "example.com/downloader/scripts/wwwinc_setup_integration"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {

	// wwwincSetupIntegration.WwwincSetupIntegration()
	// peakSetupIntegration.PeakSetupIntegration()
	// model := viewAllVersions.NewViewAllVersionsModel()
	// viewAllVersions.GetVersionsOverHttp(program)
	model := NewTopModel()
	program := tea.NewProgram(model)
	model.Program = program
	if _, err := program.Run(); err == nil {
		// fmt.Println("exiting p.Run")
	}
}
