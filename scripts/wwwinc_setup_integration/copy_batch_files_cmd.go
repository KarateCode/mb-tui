package wwwinc_setup_integration

import (
	"fmt"
)

func generateCopyFilesCmd(env wwwincEnv, batchChoice string) string {
	copyFilesCmd := fmt.Sprintf(
		`cd /client/%s/archive; cp %s.* /client/dump; cd /client/dump; ls %s.*`,
		env.subFolder,
		batchChoice,
		batchChoice,
	)

	return copyFilesCmd
}
