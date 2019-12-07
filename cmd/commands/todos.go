package commands

import (
	"context"
	"fmt"
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
		s.HideCursor = true
		s.Suffix = " finding TODOs"
		s.Writer = os.Stderr
		s.Start()

		cwd, err := os.Getwd()
		handleError(err, s)

		dir := cwd
		if len(args) == 1 {
			dir, err = filepath.Rel(cwd, args[0])
			handleError(err, s)
		}

		validateDir(dir)

		foundToDos := make(todos.ToDos, 0)
		err = comments.SearchDir(dir, func(comment *comments.Comment) {
			todo := todos.NewToDo(*comment)
			if todo != nil {
				foundToDos = append(foundToDos, todo)
				s.Suffix = fmt.Sprintf(" %d TODOs found", len(foundToDos))
			}
		})
		handleError(err, s)

		s.Suffix = fmt.Sprintf(" blaming %d TODOs", len(foundToDos))
		ctx := context.Background()
		// timeout after 30 seconds
		// ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
		// defer cancel()
		err = foundToDos.FindBlame(ctx, dir)
		sort.Sort(&foundToDos)

		handleError(err, s)

		s.Stop()

		todos.WriteTodos(foundToDos, os.Stdout)
	},
}
