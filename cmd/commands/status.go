package commands

import (
	"os"

	"github.com/augmentable-dev/tickgit"
	"github.com/spf13/cobra"
	"gopkg.in/src-d/go-git.v4"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

var versionCmd = &cobra.Command{
	Use:   "status",
	Short: "Print a status report of the current directory",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		dir, err := os.Getwd()
		handleError(err)

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
