package utils

import (
	"bufio"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
)

func IsExist(fileName string) bool {
	_, err := os.Stat(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func CleanAndExpandPath(path string) string {
	// Nothing to do when no path is given.
	if path == "" {
		return path
	}

	// NOTE: The os.ExpandEnv doesn't work with Windows cmd.exe-style
	// %VARIABLE%, but the variables can still be expanded via POSIX-style
	// $VARIABLE.
	path = os.ExpandEnv(path)

	if !strings.HasPrefix(path, "~") {
		return filepath.Clean(path)
	}

	// Expand initial ~ to the current user's home directory, or ~otheruser
	// to otheruser's home directory.  On Windows, both forward and backward
	// slashes can be used.
	path = path[1:]

	var pathSeparators string
	if runtime.GOOS == "windows" {
		pathSeparators = string(os.PathSeparator) + "/"
	} else {
		pathSeparators = string(os.PathSeparator)
	}

	userName := ""
	if i := strings.IndexAny(path, pathSeparators); i != -1 {
		userName = path[:i]
		path = path[i:]
	}

	homeDir := ""
	var u *user.User
	var err error
	if userName == "" {
		u, err = user.Current()
	} else {
		u, err = user.Lookup(userName)
	}
	if err == nil {
		homeDir = u.HomeDir
	}
	// Fallback to CWD if user lookup fails or user has no home directory.
	if homeDir == "" {
		homeDir = "."
	}

	return filepath.Join(homeDir, path)
}

func WriteLineAppendFile(file string, line string) error {
	f, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("open file error: %v\n", err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	fmt.Fprintln(w, line)
	return w.Flush()
}

func ReadLine(file string) ([]string, error) {
	f, err := os.OpenFile(file, os.O_RDONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("open file error: %v\n", err)
	}
	defer f.Close()

	rs := []string{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lineText := scanner.Text()
		rs = append(rs, lineText)

	}
	return rs, nil
}
