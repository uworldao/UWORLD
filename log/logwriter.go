package log

import (
	"github.com/jhdriver/UWORLD/log/log15/term"
	"github.com/jrick/logrotate/rotator"
	"github.com/mattn/go-colorable"
	"io"
	"os"
)

// logWriter implements an io.Writer that outputs to both standard output and
// the write-end pipe of an initialized log rotator.
type logWriter struct {
	// logRotator is one of the logging outputs.  It should be closed on
	// application shutdown.
	logRotator *rotator.Rotator

	// Use for color terminal
	colorableWrite io.Writer
}

func (lw *logWriter) Init() {
	// init a colorful logger if possible
	usecolor := term.IsTty(os.Stdout.Fd()) && os.Getenv("TERM") != "dumb"

	if usecolor {
		lw.colorableWrite = colorable.NewColorableStderr()
	}
}

func (lw *logWriter) Close() {
	if lw.logRotator != nil {
		lw.logRotator.Close()
	}
}

func (lw *logWriter) IsUseColor() bool {
	return lw.colorableWrite != nil
}

func (lw *logWriter) Write(p []byte) (n int, err error) {
	if lw.logRotator != nil {
		lw.logRotator.Write(p)
	}

	if lw.colorableWrite != nil {
		lw.colorableWrite.Write(p)
	} else {
		os.Stderr.Write(p)
	}
	return len(p), nil
}
