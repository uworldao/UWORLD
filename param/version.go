package param

import (
	"fmt"
)

var (
	GitCommitLog = "unknown_unknown"
	GitStatus    = "unknown_unknown"
	Version      = "v0.3.4"
)

func StringifySingleLine(app string) string {
	if GitStatus != "" {
		GitCommitLog = GitCommitLog[0:10] + "-dirty"
	} else {
		GitCommitLog = GitCommitLog[0:10]
	}
	return fmt.Sprintf("%s-%s-%s",
		app, Version, GitCommitLog)
}

func VersionInfo() string {
	if GitStatus != "" {
		GitCommitLog = GitCommitLog[0:10] + "-dirty"
	} else {
		GitCommitLog = GitCommitLog[0:10]
	}
	return Version + "-" + GitCommitLog
}
