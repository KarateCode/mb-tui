package peak_setup_integration

import (
	"fmt"
)

func generateCopyFilesCmd(env peakEnv, batchChoice string) string {
	prefix := calcPrefix(env.clientCode)

	copyFilesCmd := fmt.Sprintf(
		`cd /client/%s/archive; cp %s*%s.* /client/dump; cd /client/dump; ls %s*%s.*`,
		env.subFolder,
		prefix,
		batchChoice,
		prefix,
		batchChoice,
	)

	return copyFilesCmd
}
