package todos

import (
	"regexp"
	"strings"

	"github.com/augmentable-dev/tickgit/pkg/comments"
	"github.com/dustin/go-humanize"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/storer"
)

// ToDo represents a ToDo item
type ToDo struct {
	comments.Comment
	String string
	Commit *object.Commit
}

// ToDos represents a list of ToDo items
type ToDos []*ToDo

// TimeAgo returns a human readable string indicating the time since the todo was added
func (t *ToDo) TimeAgo() string {
	if t.Commit == nil {
		return "<unknown>"
	}
	return humanize.Time(t.Commit.Author.When)
}

// NewToDo produces a pointer to a ToDo from a comment
func NewToDo(comment comments.Comment) *ToDo {
	s := comment.String()
	if !strings.Contains(s, "TODO") {
		return nil
	}
	re := regexp.MustCompile(`TODO(:|,)?`)
	s = re.ReplaceAllLiteralString(comment.String(), "")
	s = strings.Trim(s, " ")

	todo := ToDo{Comment: comment, String: s}
	return &todo
}

// NewToDos produces a list of ToDos from a list of comments
func NewToDos(comments comments.Comments) ToDos {
	todos := make(ToDos, 0)
	for _, comment := range comments {
		todo := NewToDo(*comment)
		if todo != nil {
			todos = append(todos, todo)
		}
	}
	return todos
}

func (t *ToDo) existsInCommit(commit *object.Commit) (bool, error) {
	f, err := commit.File(t.FilePath)
	if err != nil {
		if err == object.ErrFileNotFound {
			return false, nil
		}
		return false, err
	}
	c, err := f.Contents()
	if err != nil {
		return false, err
	}
	return strings.Contains(c, t.Comment.String()), nil
}

// FindBlame sets the blame information on each todo in a set of todos
func (t *ToDos) FindBlame(repo *git.Repository, from *object.Commit) error {
	commitIter, err := repo.Log(&git.LogOptions{
		From: from.Hash,
	})
	if err != nil {
		return err
	}
	defer commitIter.Close()

	remainingTodos := *t
	prevCommit := from
	err = commitIter.ForEach(func(commit *object.Commit) error {
		if len(remainingTodos) == 0 {
			return storer.ErrStop
		}
		newRemainingTodos := make(ToDos, 0)
		for _, todo := range remainingTodos {
			exists, err := todo.existsInCommit(commit)
			if err != nil {
				return err
			}
			if !exists { // if the todo doesn't exist in this commit, it was added in the previous commit (previous wrt the iterator, more recent in time)
				todo.Commit = prevCommit
			} else { // if the todo does exist in this commit, add it to the new list of remaining todos
				newRemainingTodos = append(newRemainingTodos, todo)
			}
		}
		prevCommit = commit
		remainingTodos = newRemainingTodos
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
