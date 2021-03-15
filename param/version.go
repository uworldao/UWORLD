package param

import (
	"fmt"
	"runtime"
)

var (
	GitCommitLog   = "unknown_unknown"
	GitStatus      = "unknown_unknown"
	BuildTime      = "unknown_unknown"
	BuildGoVersion = "unknown_unknown"
	Version        = "v0.3.3"
)

func StringifySingleLine(app string) string {
	if GitStatus != "" {
		GitCommitLog = GitCommitLog[0:10] + "-dirty"
	} else {
		GitCommitLog = GitCommitLog[0:10]
	}
	return fmt.Sprintf("%s version=%s. commit=%s. build=%s. go=%s. runtime=%s/%s.",
		app, Version, GitStatus, BuildTime, BuildGoVersion, runtime.GOOS, runtime.GOARCH)
}

func VersionInfo() string {
	if GitStatus != "" {
		GitCommitLog = GitCommitLog[0:10] + "-dirty"
	} else {
		GitCommitLog = GitCommitLog[0:10]
	}
	return Version + "-" + GitCommitLog
}
