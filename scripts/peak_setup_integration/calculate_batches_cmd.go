package peak_setup_integration

import (
	"fmt"
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
