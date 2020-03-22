package command

import (
	"fmt"

	"github.com/NoUseFreak/go-vembed"
)

func init() {
	rootCmd.Version = fmt.Sprintf(
		"%s, build %s",
		vembed.Version.GetGitSummary(),
		vembed.Version.GetGitCommit(),
	)
}
