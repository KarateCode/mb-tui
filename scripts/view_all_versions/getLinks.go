package view_all_versions

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type ProgressGetVersion struct {
	Index         int
	VersionString string
	Error         error
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
				p.Send(ProgressGetVersion{Index: i, Error: fmt.Errorf("GET error: %w", err)})
				return
			}

			// Ensure the response body is closed to prevent resource leaks
			defer resp.Body.Close()

			// Check if the request was successful (status code 200 OK)
			if resp.StatusCode != http.StatusOK {
				p.Send(ProgressGetVersion{Index: i, Error: fmt.Errorf("status error: %d %s", resp.StatusCode, resp.Status)})
				return
			}

			// Read the entire response body into a byte slice
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				p.Send(ProgressGetVersion{Index: i, Error: fmt.Errorf("read body error: %w", err)})
				return
			}

			// Convert the byte slice to a string
			jsonData := string(bodyBytes)

			var apiResponse ApiResponse

			// Unmarshal the JSON data into the struct
			err = json.Unmarshal([]byte(jsonData), &apiResponse)
			if err != nil {
				p.Send(ProgressGetVersion{Index: i, Error: fmt.Errorf("JSON unmarshal error: %w", err)})
				return
			}

			p.Send(ProgressGetVersion{Index: i, VersionString: apiResponse.Version})
		}()
	}

	return "Download Files started"
}
