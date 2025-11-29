package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	// "os"
	// "path/filepath"
	"time"
)

func pause(s string, c chan bool) {
	fmt.Printf("Starting work on: %s\n", s)
	time.Sleep(2 * time.Second)
	fmt.Printf("Finished: %s\n", s)
	c <- true
}

func main() {
	fileNames := []string{
		"hockey_eu_product.251103012539.csv",
		"hockey_eu_pricing.251103012539.csv",
		// "hockey_eu_sku.251103012539.csv",
	}
	m := newModel(fileNames)
	p := tea.NewProgram(m)

	go DownloadFiles(fileNames, p)

	if _, err := p.Run(); err != nil {
		panic(err)
	}
	fmt.Printf("past panic check \n")

	tea.Quit()
}
