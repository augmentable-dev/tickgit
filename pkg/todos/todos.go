package todos

import (
	"bufio"
	"context"
	"regexp"
	"strings"

	"github.com/augmentable-dev/tickgit/pkg/blame"
	"github.com/augmentable-dev/tickgit/pkg/comments"
	"github.com/dustin/go-humanize"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

// ToDo represents a ToDo item
type ToDo struct {
	comments.Comment
	String string
	Blame  *blame.Blame
}

// ToDos represents a list of ToDo items
type ToDos []*ToDo

// TimeAgo returns a human readable string indicating the time since the todo was added
func (t *ToDo) TimeAgo() string {
	if t.Blame == nil {
		return "<unknown>"
	}
	return humanize.Time(t.Blame.Author.When)
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

// Len returns the number of todos
func (t ToDos) Len() int {
	return len(t)
}

// Less compares two todos by their creation time
func (t ToDos) Less(i, j int) bool {
	first := t[i]
	second := t[j]
	if first.Blame == nil || second.Blame == nil {
		return false
	}
	return first.Blame.Author.When.Before(second.Blame.Author.When)
}

// Swap swaps two todos
func (t ToDos) Swap(i, j int) {
	temp := t[i]
	t[i] = t[j]
	t[j] = temp
}

// CountWithCommits returns the number of todos with an associated commit (in which that todo was added)
func (t ToDos) CountWithCommits() (count int) {
	for _, todo := range t {
		if todo.Blame != nil {
			count++
		}
	}
	return count
}

func (t *ToDo) existsInCommit(commit *object.Commit) (bool, error) {
	f, err := commit.File(t.FilePath)
	if err != nil {
		if err == object.ErrFileNotFound {
			return false, nil
		}
		return false, err
	}
	r, err := f.Reader()
	if err != nil {
		return false, err
	}
	defer r.Close()
	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		if strings.Contains(line, t.Comment.String()) {
			return true, nil
		}
	}
	err = s.Err()
	if err != nil {
		return false, err
	}
	return false, nil
}

// FindBlame sets the blame information on each todo in a set of todos
func (t *ToDos) FindBlame(ctx context.Context, dir string) error {
	fileMap := make(map[string]ToDos)
	for _, todo := range *t {
		filePath := todo.FilePath
		if _, ok := fileMap[filePath]; !ok {
			fileMap[filePath] = make(ToDos, 0)
		}
		fileMap[filePath] = append(fileMap[filePath], todo)
	}

	for filePath, todos := range fileMap {
		lines := make([]int, 0)

		for _, todo := range todos {
			lines = append(lines, todo.StartLocation.Line)
		}
		blames, err := blame.Exec(ctx, filePath, &blame.Options{
			Directory: dir,
			Lines:     lines,
		})
		if err != nil {
			return err
		}
		for line, blame := range blames {
			for _, todo := range todos {
				if todo.StartLocation.Line == line {
					b := blame
					todo.Blame = &b
				}
			}
		}
	}
	return nil
}
