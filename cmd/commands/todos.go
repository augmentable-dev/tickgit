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
		err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				p, err := filepath.Rel(dir, path)
				handleError(err)
				_, err = todos.SearchFile(p)
				handleError(err)
				// if len(c) > 0 {
				// 	fmt.Println(c)
				// }
			}
			return nil
		})
		handleError(err)
	},
}
