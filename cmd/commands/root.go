package commands

import (
	"fmt"
	"os"

	"github.com/briandowns/spinner"
)

// TODO clean this up
func handleError(err error, spinner *spinner.Spinner) {
	if err != nil {
		if spinner != nil {
			// spinner.Suffix = ""
			spinner.FinalMSG = err.Error()
			spinner.Stop()
		} else {
			fmt.Println(err)
		}
		os.Exit(1)
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := todosCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
