package commands

import (
	"fmt"
	"os"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tickgit",
	Short: "Tickets as config",
	Long:  `tickgit is a tool for helping you manage tickets and todos in your codebase, as a part of your git history`,
}

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
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
