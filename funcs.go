package ant

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
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

// GetUserSdir returns the $HOME/.marabunta
func GetUserSdir() (string, error) {
	home := os.Getenv("HOME")
	if home == "" {
		usr, err := user.Current()
		if err != nil {
			return "", fmt.Errorf("error getting user home: %s", err)
		}
		home = usr.HomeDir
	}
	return filepath.Join(home, ".marabunta"), nil
}
