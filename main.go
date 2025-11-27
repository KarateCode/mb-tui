package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"os"
	"path/filepath"
)

func main() {
	fileName := "hockey_eu_product.251103012539.csv"
	remote := "/client/dump/" + fileName
	cwd, _ := os.Getwd()
	local := filepath.Join(cwd, fileName)

	m := newModel()
	p := tea.NewProgram(m)

	go func() {
		DownloadFile(
			"bauer-prod-eu-cf-integration",
			remote,
			local,
			func(total int64) {
				p.Send(setTotalMsg(total))
			},
			func(bytes int64) {
				p.Send(progressMsg(bytes))
			},
		)
	}()

	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
