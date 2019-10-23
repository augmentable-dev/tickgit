package commands

import (
	"os"
	"path/filepath"

	"github.com/augmentable-dev/tickgit/pkg/todos"
	"github.com/spf13/cobra"
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
		found := make([]*todos.TODO, 0)
		err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				p, err := filepath.Rel(dir, path)
				handleError(err)
				t, err := todos.SearchFile(p)
				handleError(err)
				found = append(found, t...)
			}
			return nil
		})
		handleError(err)
		todos.WriteTodos(found, os.Stdout)
	},
}
