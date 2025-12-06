package peak_setup_integration

import (
	"bufio"
	"fmt"
	"strings"
	// "time"

	exec "example.com/downloader/exec"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type batchItem string

func (i batchItem) Title() string       { return string(i) }
func (i batchItem) Description() string { return "" }
func (i batchItem) FilterValue() string { return string(i) }

type (
	downloadCompleteMsg []string
	// downloadErrorMsg error
)

type BatchModel struct {
	allBatches  []string
	filterInput textinput.Model
	list        list.Model

	isDownloading    bool
	downloadComplete bool
	showBatchesCmd   string

	Done     bool
	quitting bool
	selected string
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

func NewMenu(integrationMenuChoice IntegrationMenuChoice) BatchModel {
	choice := string(integrationMenuChoice)
	giveMeEverything := bool(choice == "Nope! Give me them all")
	requestedFileExtensions := getRequestedFileExtensions(choice)

	var showBatchesCmd string
	prefix := "hockey_eu_"
	if giveMeEverything {
		showBatchesCmd = fmt.Sprintf(
			`cd /client/EU/archive; ls | sed -n 's/%s[a-z_]*\.//p' | sed -n 's/\.csv//p' | sort | uniq | tail -n 100 | tac`,
			// One day implement env.subFolder
			prefix,
		)
	} else {
		showBatchesCmd = fmt.Sprintf(
			`cd /client/EU/archive; ls *%s* | sed -n 's/%s[a-z_]*\.//p' | sed -n 's/\.csv//p' | sort | uniq | tail -n 20 | tac`,
			requestedFileExtensions[0],
			prefix,
		)
	}
	// fmt.Printf("showBatchesCmd:\n")
	// fmt.Printf("%+v\n", showBatchesCmd)
	// tea.Quit()
	// panic("just end")

	return BatchModel{
		isDownloading:  true,
		showBatchesCmd: showBatchesCmd,
		Done:           false,
	}
}

func (m BatchModel) Init() tea.Cmd {
	return doDownload(m)
}

func doDownload(m BatchModel) tea.Cmd {
	return func() tea.Msg {
		// time.Sleep(2 * time.Second)

		host, err := exec.NewClientFromSshConfig("bauer-prod-eu-cf-integration")
		if err != nil {
			panic(err)
		}

		// prefix := "hockey_eu_"
		// cmd := fmt.Sprintf(
		// 	`cd /client/EU/archive; ls | sed -n 's/%s[a-z_]*\.//p' | sed -n 's/\.csv//p' | sort | uniq | tail -n 100 | tac`,
		// 	prefix,
		// )
		// fmt.Print("Message received!!!")
		// fmt.Print("Message received!!!")
		// fmt.Print("Message received!!!")
		// fmt.Print("Message received!!!")
		// fmt.Print("Message received!!!")
		// fmt.Print("Message received!!!")
		// fmt.Printf("m.showBatchesCmd:\n")
		// fmt.Printf("%+v\n", m.showBatchesCmd)
		// panic("just end")
		output, err := exec.RunRemoteCommand(host, m.showBatchesCmd)
		if err != nil {
			panic(err)
		}

		// fmt.Printf("output:\n")
		// fmt.Printf("%+v\n", output)
		// panic("just end")

		lines := downloadCompleteMsg{}
		// lines := []string{}
		scanner := bufio.NewScanner(strings.NewReader(string(output)))

		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line != "" {
				lines = append(lines, line)
			}
		}

		// lines := []string{"251203183900", "251203163900", "251203143900", "251203123900", "251203103900"}
		// fmt.Printf("lines:\n")
		// fmt.Printf("%+v\n", lines)

		return lines
	}
}

func setListItems(m BatchModel, lines downloadCompleteMsg) BatchModel {
	// Convert to list items
	items := make([]list.Item, len(lines))
	for i, s := range lines {
		items[i] = item(s)
	}

	// Text input
	ti := textinput.New()
	ti.Placeholder = "Search..."
	ti.Focus()

	// List
	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = false
	delegate.SetHeight(1)
	delegate.SetSpacing(0)

	l := list.New(items, delegate, 50, 40) // WIDTH=50, HEIGHT=20 rows
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.SetShowHelp(false)
	l.SetShowTitle(false)
	l.SetShowPagination(false)

	l.Styles.Title = lipgloss.NewStyle()

	m.allBatches = lines
	m.filterInput = ti
	m.list = l
	return m
}

func (m BatchModel) Update(msg tea.Msg) (BatchModel, tea.Cmd) {
	switch msg := msg.(type) {

	case downloadCompleteMsg:
		batchModel := setListItems(m, msg)
		batchModel.isDownloading = false
		batchModel.downloadComplete = true
		return batchModel, nil
		// m.isDownloading = false
		// m.downloadComplete = true
		// return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit

		// Emacs-style movement
		case "ctrl+n":
			// Move down
			m.list.CursorDown()
			return m, nil

		case "ctrl+p":
			// Move up
			m.list.CursorUp()
			return m, nil

		case "enter":
			selected := fmt.Sprintf("%+v", m.list.SelectedItem())
			// selected, ok := m.list.SelectedItem().(batchItem)
			// if ok {
			m.selected = string(selected)
			m.Done = true
			// }
			// if selected, ok := m.list.SelectedItem().(batchItem); ok {
			// 	m.selected = string(selected)
			// 	m.Done = true
			// }
			// fmt.Printf("selected:\n")
			// fmt.Printf("selected:\n")
			// fmt.Printf("selected:\n")
			// fmt.Printf("selected:\n")
			// fmt.Printf("selected:\n")
			// fmt.Printf("selected:\n")
			// fmt.Printf("m.list:\n")
			// fmt.Printf("%+v\n", m.list)
			// fmt.Printf("selected:\n")
			// fmt.Printf("%+v\n", selected)
			// fmt.Printf("ok:\n")
			// fmt.Printf("%+v\n", ok)
			// panic("just end")
			return m, nil
			// return m, tea.Quit
		}
	}

	// Update text input
	var cmd tea.Cmd
	m.filterInput, cmd = m.filterInput.Update(msg)

	// Filter list items
	filter := strings.ToLower(m.filterInput.Value())

	// var filtered []list.Item
	// for _, b := range m.allBatches {
	// 	if strings.Contains(strings.ToLower(b), filter) {
	// 		filtered = append(filtered, item(b))
	// 	}
	// }
	// m.list.SetItems(filtered)
	if filter != "" {
		var filtered []list.Item
		for _, b := range m.allBatches {
			if strings.Contains(strings.ToLower(b), filter) {
				filtered = append(filtered, batchItem(b))
			}
		}
		m.list.SetItems(filtered)
	} else {
		// Only reset when returning to full list
		m.list.SetItems(batchItemsFrom(m.allBatches))
	}

	// Update list
	m.list, _ = m.list.Update(msg)

	return m, cmd
}

// itemsFrom converts []string -> []list.Item
func batchItemsFrom(batches []string) []list.Item {
	items := make([]list.Item, 0, len(batches))
	for _, b := range batches {
		items = append(items, item(b))
	}
	return items
}

func (m BatchModel) View() string {
	if m.quitting {
		return ""
	}

	if m.isDownloading {
		return "Downloading file..."
		// return fmt.Sprintf("%s Downloading file...", m.spinner.View())
	}

	if m.downloadComplete {
		// return "✅ Download complete! Press Ctrl+C to exit.\n"
		return fmt.Sprintf(
			"✅ Download complete! Press Ctrl+C to exit.\nFilter: %s\n\n%s",
			m.filterInput.View(),
			m.list.View(),
		)
	}

	return "unknown state from batchmenu.model"
}

func (m BatchModel) Selected() string {
	return m.selected
}
