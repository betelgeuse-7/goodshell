package main

import (
	"fmt"
	"os"
)

func cd(argv []string) error {
	var goToHome = func() error {
		// UserHomeDir returns an error only when $HOME is not defined;
		// and that possibility is very slim.
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("Unable to get $HOME\n")
		}
		if err := os.Chdir(home); err != nil {
			return err
		}
		return nil
	}
	if len(argv) > 1 {
		return fmt.Errorf("cd: excessive arguments\n")
	}
	if len(argv) == 0 {
		return goToHome()
	}
	dir := argv[0]
	if dir == "~" {
		return goToHome()
	}
	err := os.Chdir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("cd: No such directory: %s\n", dir)
		}
		return fmt.Errorf("This should not be printed. This is a bug\n%s\n", err)
	}
	return nil
}
