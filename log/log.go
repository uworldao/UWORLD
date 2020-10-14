package log

import (
	"fmt"
	log "github.com/jhdriver/UWORLD/log/log15"
	"github.com/jrick/logrotate/rotator"
	"os"
	"path/filepath"
)

var (
	glogger  *log.GlogHandler
	logWrite *logWriter
)

func init() {
	logWrite = &logWriter{}
	logWrite.Init()
	glogger = log.NewGlogHandler(log.StreamHandler(logWrite, log.TerminalFormat(true)))
	log.Root().SetHandler(glogger)

	glogger.Verbosity(log.LvlInfo)
}

func InitLogRotator(logFile string) {
	logDir, _ := filepath.Split(logFile)
	err := os.MkdirAll(logDir, 0700)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create log directory: %v\n", err)
		os.Exit(1)
	}
	r, err := rotator.New(logFile, 10*1024, false, 3)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create file rotator: %v\n", err)
		os.Exit(1)
	}

	logWrite.logRotator = r
}

func LogWrite() *logWriter {
	return logWrite
}

func Glogger() *log.GlogHandler {
	return glogger
}
