package peak_setup_integration

type peakEnv struct {
	name       string
	sshServer  string
	subFolder  string
	clientCode string
}
type peakEnvs []peakEnv

func integrationListItems() []string {
	return []string{
		"Nope! Give me them all",
		"Product Import",
		"Customer Import",
		"Inventory Import",
		"SalesRep Import",
		"BG/BHC import",
		"SalesOrg/PoType Import",
	}
}

func environments() peakEnvs {
	return peakEnvs{
		{
			name:       "Bauer EU Staging",
			sshServer:  "bauer-stag-eu-cf-integration",
			subFolder:  "EU",
			clientCode: "bauer-eu",
		},
		{
			name:       "Bauer EU Production",
			sshServer:  "bauer-prod-eu-cf-integration",
			subFolder:  "EU",
			clientCode: "bauer-eu",
		},
		{
			name:       "Bauer NA Staging",
			sshServer:  "bauer-stag-na-cf-integration",
			subFolder:  "NA",
			clientCode: "bauer-na",
		},
		{
			name:       "Bauer NA Production",
			sshServer:  "bauer-prod-na-cf-integration",
			subFolder:  "NA",
			clientCode: "bauer-na",
		},
		{
			name:       "Cascade NA Staging",
			sshServer:  "cascade-stag-na-cf-integration",
			subFolder:  "NA",
			clientCode: "cascade-na",
		},
		{
			name:       "Cascade NA Production",
			sshServer:  "cascade-prod-na-cf-integration",
			subFolder:  "NA",
			clientCode: "cascade-na",
		},
	}
}

// Define a function that mimics _.find behavior
func findEnvByName(envs peakEnvs, selectedName string) (peakEnv, bool) {
	for _, env := range envs {
		if env.name == selectedName {
			return env, true
		}
	}

	return peakEnv{}, false
}
