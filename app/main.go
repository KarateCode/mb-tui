package main

import (
	scripts "example.com/downloader/scripts"
)

func main() {
	scripts.PeakSetupIntegration()
	// host, err := exec.NewClientFromSshConfig("bauer-prod-eu-cf-integration")
	// if err != nil {
	// 	panic(err)
	// }

	// prefix := "hockey_eu_"
	// cmd := fmt.Sprintf(
	// 	`cd /client/EU/archive; ls | sed -n 's/%s[a-z_]*\.//p' | sed -n 's/\.csv//p' | sort | uniq | tail -n 100 | tac`,
	// 	prefix,
	// )
	// output, err := exec.RunRemoteCommand(host, cmd)
	// if err != nil {
	// 	panic(err)
	// }

	// lines := []string{}
	// scanner := bufio.NewScanner(strings.NewReader(string(output)))

	// for scanner.Scan() {
	// 	line := strings.TrimSpace(scanner.Text())
	// 	if line != "" {
	// 		lines = append(lines, line)
	// 	}
	// }
	// fmt.Printf("lines:\n")
	// fmt.Printf("%+v\n", lines)

	// m := batchmenu.NewMenu(lines)
	// p := tea.NewProgram(m)
	// if _, err := p.Run(); err == nil {
	// 	// fmt.Println("Selected:", result.(batchmenu.Model).Selected())
	// 	fmt.Println("exiting p.Run")
	// }

	// return

	// fileNames := []string{
	// 	"hockey_eu_product.251103012539.csv",
	// 	"hockey_eu_pricing.251103012539.csv",
	// 	"hockey_eu_sku.251103012539.csv",
	// }
	// m := newModel(fileNames)
	// p := tea.NewProgram(m)

	// go DownloadFiles(fileNames, p)

	// if _, err := p.Run(); err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("\nDownload complete! 🎉 \n")

	// tea.Quit()
}
