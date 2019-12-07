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

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Print a status report of the current directory",
	Long:  ``,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cwd, err := os.Getwd()
		handleError(err, nil)

		dir := cwd
		if len(args) == 1 {
			dir, err = filepath.Rel(cwd, args[0])
			handleError(err, nil)
		}

		validateDir(dir)

		r, err := git.PlainOpen(dir)
		handleError(err, nil)

		ref, err := r.Head()
		handleError(err, nil)

		commit, err := r.CommitObject(ref.Hash())
		handleError(err, nil)

		err = tickgit.WriteStatus(commit, os.Stdout)
		handleError(err, nil)
	},
}
