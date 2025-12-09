package peak_setup_integration

import (
	// "bufio"
	"fmt"
	// "strings"

	exec "example.com/downloader/exec"
	tea "github.com/charmbracelet/bubbletea"
)

type calculateBatchesModel struct {
	showBatchesCmd string
	peakEnv        peakEnv
}

func getRequestedFileExtensions(choice string) []string {
	if choice == "Inventory Import" {
		return []string{"inventory"}
	} else if choice == "BG/BHC import" {
		return []string{"bg_bhc"}
	} else if choice == "SalesOrg/PoType Import" {
		return []string{"salesorg_po_type"}
	} else if choice == "Customer Import" {
		return []string{"customer"}
	} else if choice == "Product Import" {
		return []string{
			"product",
			"sku",
			"pricing",
		}
	} else if choice == "SalesRep Import" {
		return []string{"salesrep"}
	}

	return nil
}

func commandForIntegration(choice string, env peakEnv) string {
	giveMeEverything := bool(choice == "Nope! Give me them all")
	requestedFileExtensions := getRequestedFileExtensions(choice)

	prefix := calcPrefix(env.clientCode)
	var showBatchesCmd string
	if giveMeEverything {
		showBatchesCmd = fmt.Sprintf(
			`cd /client/%s/archive; ls | sed -n 's/%s[a-z_]*\.//p' | sed -n 's/\.csv//p' | sort | uniq | tail -n 100 | tac`,
			env.subFolder,
			prefix,
		)
	} else {
		showBatchesCmd = fmt.Sprintf(
			`cd /client/%s/archive; ls *%s* | sed -n 's/%s[a-z_]*\.//p' | sed -n 's/\.csv//p' | sort | uniq | tail -n 20 | tac`,
			env.subFolder,
			requestedFileExtensions[0],
			prefix,
		)
	}
	return showBatchesCmd
}

func NewCalculateBatchesModel(integrationMenuChoice IntegrationMenuChoice, env peakEnv) calculateBatchesModel {
	choice := string(integrationMenuChoice)
	giveMeEverything := bool(choice == "Nope! Give me them all")
	requestedFileExtensions := getRequestedFileExtensions(choice)

	prefix := calcPrefix(env.clientCode)
	var showBatchesCmd string
	if giveMeEverything {
		showBatchesCmd = fmt.Sprintf(
			`cd /client/%s/archive; ls | sed -n 's/%s[a-z_]*\.//p' | sed -n 's/\.csv//p' | sort | uniq | tail -n 100 | tac`,
			env.subFolder,
			prefix,
		)
	} else {
		showBatchesCmd = fmt.Sprintf(
			`cd /client/%s/archive; ls *%s* | sed -n 's/%s[a-z_]*\.//p' | sed -n 's/\.csv//p' | sort | uniq | tail -n 20 | tac`,
			env.subFolder,
			requestedFileExtensions[0],
			prefix,
		)
	}

	return calculateBatchesModel{
		showBatchesCmd: showBatchesCmd,
		peakEnv:        env,
	}
}

func (m calculateBatchesModel) Init() tea.Cmd {
	return func() tea.Msg {
		host, err := exec.NewClientFromSshConfig(m.peakEnv.sshServer)
		if err != nil {
			panic(err)
		}

		output, err := exec.RunRemoteCommand(host, m.showBatchesCmd)
		if err != nil {
			panic(err)
		}

		return calcBatchesCompleteMsg(output)
		// lines := calcBatchesCompleteMsg{}
		// scanner := bufio.NewScanner(strings.NewReader(string(output)))

		// for scanner.Scan() {
		// 	line := strings.TrimSpace(scanner.Text())
		// 	if line != "" {
		// 		lines = append(lines, line)
		// 	}
		// }

		// return lines
	}
}

func (m calculateBatchesModel) Update(msg tea.Msg) (calculateBatchesModel, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m calculateBatchesModel) View() string {
	return "\n\nCalculating batches..."
}
