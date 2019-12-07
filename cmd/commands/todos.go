package commands

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/augmentable-dev/tickgit/pkg/comments"
	"github.com/augmentable-dev/tickgit/pkg/todos"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(todosCmd)
}

var todosCmd = &cobra.Command{
	Use:   "todos",
	Short: "Print a report of current TODOs",
	Long:  `Scans a given git repository looking for any code comments with TODOs. Displays a report of all the TODO items found.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Suffix = " finding TODOs"
		s.Writer = os.Stderr
		s.Start()

		cwd, err := os.Getwd()
		handleError(err)

		dir := cwd
		if len(args) == 1 {
			dir, err = filepath.Rel(cwd, args[0])
			handleError(err)
		}

		validateDir(dir)

		comments, err := comments.SearchDir(dir)
		handleError(err)

		t := todos.NewToDos(comments)

		ctx := context.Background()
		// timeout after 30 seconds
		// ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
		// defer cancel()
		err = t.FindBlame(ctx, dir)
		sort.Sort(&t)

		handleError(err)

		s.Stop()
		todos.WriteTodos(t, os.Stdout)
	},
}
