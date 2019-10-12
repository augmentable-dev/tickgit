package tickgit

import (
	"math"
	"path/filepath"

	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

// DefaultMatchPatterns are the default path/filepath match pattern strings used to detect tickgit files
var DefaultMatchPatterns []string = []string{"*.tickgit"}

// File represents a tickgit file
type File struct {
	Goals []Goal `hcl:"goal,block"`
}

// Goal represents a goal, which houses a collection of tasks
type Goal struct {
	Title string  `hcl:"title,label"`
	Tasks []*Task `hcl:"task,block"`
}

// Task represents a task
type Task struct {
	Title       string  `hcl:"title,label"`
	Description *string `hcl:"description"`
	Status      string  `hcl:"status"`
}

// TaskSummary is a summary of a set of tasks
type TaskSummary struct {
	Total     int
	Pending   int
	Completed int
}

// GoalSummary is a summary of a set of goals
type GoalSummary struct {
	Total     int
	Pending   int
	Completed int
}

// Parse parses an HCL config string
func Parse(src []byte, filename string) (*File, error) {
	file := &File{}
	parser := hclparse.NewParser()
	f, diag := parser.ParseHCL(src, filename)
	if diag != nil {
		return nil, diag
	}
	diag = gohcl.DecodeBody(f.Body, nil, file)
	if diag != nil {
		return nil, diag
	}
	return file, nil
}

// Completed returns whether or not the task can be considered completed
func (task *Task) Completed() bool {
	if task.Status == "done" {
		return true
	}
	return false
}

// Summary summarizes the tasks in a goal
func (goal *Goal) Summary() *TaskSummary {
	summary := &TaskSummary{Total: len(goal.Tasks)}
	for _, task := range goal.Tasks {
		if task.Completed() {
			summary.Completed++
		} else {
			summary.Pending++
		}
	}
	return summary
}

// TasksCompleted returns whether all the tasks in this summary are completed
func (taskSummary *TaskSummary) TasksCompleted() bool {
	return taskSummary.Completed == taskSummary.Total
}

// Completed returns whether all the tasks of a goal are completed
func (goal *Goal) Completed() bool {
	summary := goal.Summary()
	return summary.TasksCompleted()
}

// PercentCompleted returns the percentage of tasks completed
func (taskSummary *TaskSummary) PercentCompleted() int {
	f := math.Round(float64(taskSummary.Completed) / float64(taskSummary.Total) * 100)
	return int(f)
}

// GoalsFromCommit returns all the goals in the given commit's tree
func GoalsFromCommit(commit *object.Commit, matchPatternOverrides []string) ([]Goal, error) {
	matchPatterns := DefaultMatchPatterns

	// if the caller supplied matchPatternOverrides, use it
	if matchPatternOverrides != nil || len(matchPatternOverrides) != 0 {
		matchPatterns = matchPatternOverrides
	}
	goals := make([]Goal, 0)

	// get the tree of the commit we want to retrieve goals from
	tree, err := commit.Tree()
	if err != nil {
		return nil, err
	}

	// search every file in the tree, looking for a match
	err = tree.Files().ForEach(func(file *object.File) error {
		matched := false
		for _, matchPattern := range matchPatterns {
			currentMatch, err := filepath.Match(matchPattern, filepath.Base(file.Name))
			if err != nil {
				return err
			}

			if currentMatch {
				matched = true
				break
			}
		}

		if matched {
			content, err := file.Contents()
			if err != nil {
				return err
			}
			f, err := Parse([]byte(content), file.Name)
			if err != nil {
				return err
			}
			goals = append(goals, f.Goals...)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return goals, nil
}
