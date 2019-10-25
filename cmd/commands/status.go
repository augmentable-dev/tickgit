package commands

import (
	"os"
	"path/filepath"

	"github.com/augmentable-dev/tickgit"
	"github.com/spf13/cobra"
	"gopkg.in/src-d/go-git.v4"
)

func init() {
	rootCmd.AddCommand(statusCmd)
}

// TODO clean this up
func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Print a status report of the current directory",
	Long:  ``,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cwd, err := os.Getwd()
		handleError(err)

		dir := cwd
		if len(args) == 1 {
			dir, err = filepath.Rel(cwd, args[0])
			handleError(err)
		}

		r, err := git.PlainOpen(dir)
		handleError(err)

		ref, err := r.Head()
		handleError(err)

		commit, err := r.CommitObject(ref.Hash())
		handleError(err)

		err = tickgit.WriteStatus(commit, os.Stdout)
		handleError(err)
	},
}
