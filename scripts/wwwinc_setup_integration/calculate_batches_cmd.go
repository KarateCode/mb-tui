package wwwinc_setup_integration

import (
	"fmt"
)

func getRequestedFileExtensions(choice string) []string {
	if choice == "User import" {
		return []string{
			"CUSTOMERSECURITY",
			"SALESREPSECURITY",
			"CSR",
			"CUS",
			"SRM",
			"DON",
		}
	} else if choice == "Catalog import" {
		return []string{
			"CATALOGLANG",
			"CATALOG",
			"CATALOGITM",
			"CAT",
			"CATITM",
			"CUSTOMERSECURITY",
			"BRANDMASTER",
			"FLC",
			"FLG",
			"DON",
		}
	} else if choice == "Customer import" {
		return []string{
			"CATALOG",
			"CATALOGCUS",
			"CSH",
			"BRANDMASTER",
			"CUSTOMERSECURITY",
			"CUS",
			"SALESREPSECURITY",
			"ITE",
			"DON",
		}
	} else if choice == "Flag import" {
		return []string{
			"FLG",
			"BRANDMASTER",
			"DON",
		}
	} else if choice == "Grid import" {
		return []string{
			"ITE",
			"BRANDMASTER",
			"MGR",
			"DON",
		}
	} else if choice == "Product import" {
		return []string{
			"ITE",
			"MKT",
			"FSA",
			"PRICE",
			"BRANDMASTER",
			"UPC",
			"DON",
		}
	} else if choice == "Sales Rep import" {
		return []string{
			"SALESREPSECURITY",
			"BRANDMASTER",
			"SRM",
			"DON",
		}
	}

	return nil
}

func commandForIntegration(choice string, env wwwincEnv) string {
	giveMeEverything := bool(choice == "Nope! Give me them all")
	requestedFileExtensions := getRequestedFileExtensions(choice)

	var showBatchesCmd string
	if giveMeEverything {
		requestedFileExtension := "CATALOGLANG"
		showBatchesCmd = fmt.Sprintf(
			`cd /client/%s/archive; ls *.%s | sed -n 's/\..*//p' | sort | uniq | tail -n 20 | tac`,
			env.subFolder,
			requestedFileExtension,
		)
	} else {
		showBatchesCmd = fmt.Sprintf(
			`cd /client/%s/archive; ls *.%s | sed -n 's/\..*//p' | sort | uniq | tail -n 20 | tac`,
			env.subFolder,
			requestedFileExtensions[0],
		)
	}
	return showBatchesCmd
}
