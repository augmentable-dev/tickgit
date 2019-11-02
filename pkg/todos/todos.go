package todos

import (
	"regexp"
	"strings"

	"github.com/augmentable-dev/tickgit/pkg/comments"
)

// ToDo represents a ToDo item
type ToDo struct {
	comments.Comment
	String string
}

// ToDos represents a list of ToDo items
type ToDos []ToDo

// Count returns the number of todos
func (t ToDos) Count() int {
	return len(t)
}

// NewToDo produces a pointer to a ToDo from a comment
func NewToDo(comment comments.Comment) *ToDo {
	s := comment.String()
	if !strings.Contains(s, "TODO") {
		return nil
	}
	re := regexp.MustCompile(`TODO:?`)
	s = re.ReplaceAllLiteralString(comment.String(), "")
	s = strings.Trim(s, " ")

	todo := ToDo{Comment: comment, String: s}
	return &todo
}

// NewToDos produces a list of ToDos from a list of comments
func NewToDos(comments comments.Comments) ToDos {
	todos := make(ToDos, 0)
	for _, comment := range comments {
		todo := NewToDo(comment)
		if todo != nil {
			todos = append(todos, *todo)
		}
	}
	return todos
}
