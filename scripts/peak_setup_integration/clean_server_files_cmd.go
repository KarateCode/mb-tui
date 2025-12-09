package peak_setup_integration

import (
	"fmt"
)

func generateCleanServerFilesCmd(env peakEnv, batchChoice string) string {
	prefix := calcPrefix(env.clientCode)

	cmd := fmt.Sprintf(
		"cd /client/dump; rm %s*%s.*",
		prefix,
		batchChoice,
	)

	return cmd
}
