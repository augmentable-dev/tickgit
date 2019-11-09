package todos

import (
	"regexp"
	"strings"
	"time"

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
	Added  *time.Time
	Author string
}

// ToDos represents a list of ToDo items
type ToDos []*ToDo

// Count returns the number of todos
func (t ToDos) Count() int {
	return len(t)
}

// TimeAgo returns a human readable string indicating the time since the todo was added
func (t *ToDo) TimeAgo() string {
	if t.Added == nil {
		return "<unknown>"
	}
	return humanize.Time(*t.Added)
	// dur := time.Now().Sub(*t.Added)

	// hours := dur.Hours()
	// days := hours / 24
	// weeks := days / 7
	// months := days / 30
	// years := months / 12

	// if hours <= 24 {
	// 	return fmt.Sprintf("~%d hours ago", int(math.Round(hours)))
	// } else if days <= 7 {
	// 	return fmt.Sprintf("~%d days ago", int(math.Round(days)))
	// } else if weeks <= 4 {
	// 	return fmt.Sprintf("~%d weeks ago", int(math.Round(weeks)))
	// } else if months <= 12 {
	// 	return fmt.Sprintf("~%d months ago", int(math.Round(months)))
	// } else {
	// 	return fmt.Sprintf("~%d years ago", int(math.Round(years)))
	// }
}

// FindBlame sets the Added and Author fields on the ToDo
// TODO: find ways to optimize this, set a ceiling to stop searching the history after some time
// run this against a batch of todos so that git history is not traversed per-todo
func (t *ToDo) FindBlame(repo *git.Repository, from *object.Commit) error {
	commitIter, err := repo.Log(&git.LogOptions{
		From: from.Hash,
	})
	if err != nil {
		return err
	}
	defer commitIter.Close()

	prevCommit := from
	err = commitIter.ForEach(func(commit *object.Commit) error {
		prevCommit = commit
		f, err := commit.File(t.FilePath)
		if err != nil {
			return err
		}
		c, err := f.Contents()
		if err != nil {
			return err
		}
		todoPresent := strings.Contains(c, t.String)

		if !todoPresent {
			return storer.ErrStop
		}
		return nil
	})
	t.Author = prevCommit.Author.String()
	t.Added = &(prevCommit.Author.When)
	if err != nil {
		return nil
	}
	return nil
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
