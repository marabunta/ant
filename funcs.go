package ant

import (
	"os"
)

// isDir return true if path is a dir
func isDir(path string) bool {
	f, err := os.Stat(path)
	if err != nil {
		return false
	}
	if m := f.Mode(); m.IsDir() && m&400 != 0 {
		return true
	}
	return false
}

// isFile return true if path is a regular file
func isFile(path string) bool {
	f, err := os.Stat(path)
	if err != nil {
		return false
	}
	if m := f.Mode(); !m.IsDir() && m.IsRegular() && m&400 != 0 {
		return true
	}
	return false
}
