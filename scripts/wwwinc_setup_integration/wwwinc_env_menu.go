package wwwinc_setup_integration

type wwwincEnv struct {
	name       string
	sshServer  string
	subFolder  string
	clientCode string
}
type wwwincEnvs []wwwincEnv

func integrationListItems() []string {
	return []string{
		"Nope! Give me them all",
		"Catalog import",
		"Customer import",
		"Flag import",
		"Grid import",
		"Sales Rep import",
		"Product import",
		"User import",
	}
}

func environments() wwwincEnvs {
	return wwwincEnvs{
		{
			name:       "Wwwinc NA Staging",
			sshServer:  "wwwinc-stag-na-cf-integration",
			subFolder:  "US",
			clientCode: "wwwinc-us",
		},
		{
			name:       "Wwwinc NA Production",
			sshServer:  "wwwinc-prod-na-cf-integration",
			subFolder:  "US",
			clientCode: "wwwinc-us",
		},
		{
			name:       "Wwwinc EU Staging",
			sshServer:  "wwwinc-stag-eu-cf-integration",
			subFolder:  "EMEA",
			clientCode: "wwwinc-emea",
		},
		{
			name:       "Wwwinc EU Production",
			sshServer:  "wwwinc-prod-eu-cf-integration",
			subFolder:  "EU",
			clientCode: "wwwinc-emea",
		},
	}
}

// Define a function that mimics _.find behavior
func findEnvByName(envs wwwincEnvs, selectedName string) (wwwincEnv, bool) {
	for _, env := range envs {
		if env.name == selectedName {
			return env, true
		}
	}

	return wwwincEnv{}, false
}
