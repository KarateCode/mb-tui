package wwwinc_setup_integration

import (
	"fmt"
)

func generateCleanServerFilesCmd(batchChoice string) string {
	cmd := fmt.Sprintf(
		"cd /client/dump; rm %s.*",
		batchChoice,
	)

	return cmd
}
