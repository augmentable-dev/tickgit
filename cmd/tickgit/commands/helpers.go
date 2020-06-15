package commands

import (
	"fmt"
	"os"
	"path/filepath"
)

func validateDir(dir string) {
	if dir == "" {
		cwd, err := os.Getwd()
		handleError(err, nil)
		dir = cwd
	}

	abs, err := filepath.Abs(filepath.Join(dir, ".git"))
	handleError(err, nil)

	if _, err := os.Stat(abs); os.IsNotExist(err) {
		handleError(fmt.Errorf("%s is not a git repository", abs), nil)
	}
}
