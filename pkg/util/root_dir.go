package util

import (
	"path/filepath"
	"runtime"
)

// RootDir is a function uses to get an absolute path of current package.
func RootDir() string {
	_, b, _, _ := runtime.Caller(0)

	return filepath.Join(filepath.Dir(b), "..", "..")
}
