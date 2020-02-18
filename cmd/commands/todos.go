package commands

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"

	"github.com/augmentable-dev/tickgit/pkg/comments"
	"github.com/augmentable-dev/tickgit/pkg/todos"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
)

var csvOutput bool

func init() {
	todosCmd.Flags().BoolVar(&csvOutput, "csv-output", false, "specify whether or not output should be in CSV format")
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

		if csvOutput {
			w := csv.NewWriter(os.Stdout)
			err := w.Write([]string{
				"text", "file_path", "start_line", "start_position", "end_line", "end_position", "author", "author_email", "author_sha", "author_time",
			})
			handleError(err, s)

			for _, todo := range foundToDos {
				err := w.Write([]string{
					todo.String,
					todo.FilePath,
					strconv.Itoa(todo.StartLocation.Line),
					strconv.Itoa(todo.StartLocation.Pos),
					strconv.Itoa(todo.EndLocation.Line),
					strconv.Itoa(todo.EndLocation.Pos),
					todo.Blame.Author.Name,
					todo.Blame.Author.Email,
					todo.Blame.SHA,
					todo.Blame.Author.When.Format(time.RFC3339),
				})
				handleError(err, s)
			}

			// Write any buffered data to the underlying writer (standard output).
			w.Flush()

			err = w.Error()
			handleError(err, s)

		} else {
			todos.WriteTodos(foundToDos, os.Stdout)
		}

	},
}
