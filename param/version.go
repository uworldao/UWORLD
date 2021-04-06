package param

import (
	"fmt"
)

var (
	GitCommitLog = "unknown_unknown"
	GitStatus    = "unknown_unknown"
	Version      = "v0.3.4"
)

func StringifySingleLine() string {
	if len(GitCommitLog) < 7 {
		return ""
	}
	if GitStatus != "" {
		GitCommitLog = GitCommitLog[0:7] + "-dirty"
	} else {
		GitCommitLog = GitCommitLog[0:7]
	}
	return fmt.Sprintf("%s-%s-%s",
		AppName, Version, GitCommitLog)
}
