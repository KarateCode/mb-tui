package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	fileNames := []string{
		"hockey_eu_product.251103012539.csv",
		"hockey_eu_pricing.251103012539.csv",
		"hockey_eu_sku.251103012539.csv",
	}
	m := newModel(fileNames)
	p := tea.NewProgram(m)

	go DownloadFiles(fileNames, p)

	if _, err := p.Run(); err != nil {
		panic(err)
	}
	fmt.Printf("\nDownload complete! 🎉 \n")

	tea.Quit()
}
