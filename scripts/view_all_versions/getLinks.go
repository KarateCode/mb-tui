package view_all_versions

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type ProgressGetVersion struct {
	Index         int
	VersionString string
}

type ApiResponse struct {
	Version string    `json:"version"`
	Build   string    `json:"build"`
	Time    time.Time `json:"time"`
}

var links = []string{
	"https://wwwinc.na.envoy-staging.com/version.json",
	"https://login.orderwwwbrands.com/version.json",
	"https://wwwinc.eu.envoy-staging.com/version.json",
	"https://login.stag.bauer.eu.envoyb2bstaging.com/version.json",
	"https://login.stag.bauer.na.envoyb2bstaging.com/version.json",
	"https://login.bauerb2b.com/version.json",
	"https://login.stag.bauer.eu.envoyb2bstaging.com/version.json",
	"https://login.bauerb2b.hockey/version.json",
	"https://login.stag.cascademaverik.na.envoyb2bstaging.com/version.json",
	"https://login.cascademaverikb2b.com/version.json",
	"https://converse.na.envoy-staging.com/version.json",
	"https://us.converse.net/version.json",
	"https://converse.eu.envoy-staging.com/version.json",
	"https://eu.converse.net/version.json",
	"https://converse.la.envoy-staging.com/version.json",
	"https://la.converse.net/version.json",
	"https://converse.ap.envoy-staging.com/version.json",
	"https://ap.converse.net/version.json",
	"https://envoylogin.envoyb2bstaging.com/version.json",
	"https://login.envoyb2b.com/version.json",
	"https://tank-na-ns.ns.envoyb2bstaging.com/version.json",
	"https://hc-mfg.envoyb2b.com/version.json",

	"https://login.envoydemo.com/version.json",

	"https://cid.cid-resources.na.envoyb2bstaging.com/version.json",
	"https://cidone.cidresources.com/version.json",
	"https://stag.danpost.envoyb2bstaging.com/version.json",
	"https://danpostb2b.com/version.json",
	"https://oofos.na.envoy-staging.com/version.json",
	"https://oofos.oofosb2b.com/version.json",
	"https://vidagroup.na.envoy-staging.com/version.json",
	"https://ordervidabrands.com/version.json",
	"https://landau.na.envoy-staging.com/version.json",
	"https://b2b.landau.com/version.json",
	"https://uk.joules.eu.envoy-staging.com/version.json",
	"https://uk.joules.eu.envoy-staging.com/version.json",
}

func GetVersionsOverHttp(p *tea.Program) tea.Msg {
	for i, link := range links {

		go func() {
			resp, err := http.Get(link)
			if err != nil {
				log.Fatal("GET error:", err)

			}

			// Ensure the response body is closed to prevent resource leaks
			defer resp.Body.Close()

			// Check if the request was successful (status code 200 OK)
			if resp.StatusCode != http.StatusOK {
				log.Fatal("Status error:", resp.StatusCode, resp.Status)

			}

			// Read the entire response body into a byte slice
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Fatal("Read body error:", err)

			}

			// Convert the byte slice to a string
			jsonData := string(bodyBytes)

			var apiResponse ApiResponse

			// Unmarshal the JSON data into the struct
			err = json.Unmarshal([]byte(jsonData), &apiResponse)
			if err != nil {
				log.Fatalf("Error unmarshalling JSON: %v", err)

			}

			p.Send(ProgressGetVersion{Index: i, VersionString: apiResponse.Version})
		}()
	}

	return "Download Files started"
}
