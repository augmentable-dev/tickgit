package commands

import (
	"os"
	"path/filepath"

	"github.com/augmentable-dev/tickgit/pkg/comments"
	"github.com/augmentable-dev/tickgit/pkg/todos"
	"github.com/spf13/cobra"
	"gopkg.in/src-d/go-git.v4"
)

func init() {
	rootCmd.AddCommand(todosCmd)
}

var todosCmd = &cobra.Command{
	Use:   "todos",
	Short: "Print a report of current TODOs",
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

		comments, err := comments.SearchCommit(commit)
		handleError(err)

		// comments, err := comments.SearchDir(dir)
		// handleError(err)

		t := todos.NewToDos(comments)
		todos.WriteTodos(t, os.Stdout)
	},
}
