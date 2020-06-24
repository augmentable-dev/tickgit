package todos

import (
	"context"
	"strings"

	"github.com/augmentable-dev/tickgit/pkg/blame"
	"github.com/augmentable-dev/tickgit/pkg/comments"
	"github.com/dustin/go-humanize"
)

// ToDo represents a ToDo item
type ToDo struct {
	comments.Comment
	String string
	Phrase string
	Blame  *blame.Blame
}

// ToDos represents a list of ToDo items
type ToDos []*ToDo

var startingMatchPhrases []string = []string{"TODO", "FIXME", "OPTIMIZE", "HACK", "XXX", "WTF", "LEGACY"}

// TimeAgo returns a human readable string indicating the time since the todo was added
func (t *ToDo) TimeAgo() string {
	if t.Blame == nil {
		return "<unknown>"
	}
	return humanize.Time(t.Blame.Author.When)
}

// NewToDo produces a pointer to a ToDo from a comment
func NewToDo(comment comments.Comment, matchPhrases []string) *ToDo {
	// FIXME this should be configurable and probably NOT hardcoded here
	// in fact, this list might be too expansive for a sensible defaul
	for _, phrase := range matchPhrases {
		// populates matchPhrases with the contents of startingMatchPhrases plus the @+lowerCase version of each phrase
		matchPhrases = append(matchPhrases, phrase, "@"+strings.ToLower(phrase))
	}

	for _, phrase := range matchPhrases {
		s := comment.String()
		if strings.Contains(s, phrase) {
			todo := ToDo{
				Comment: comment,
				String:  strings.Trim(s, " "),
				Phrase:  phrase,
			}
			return &todo
		}
	}

	return nil
}

// NewToDos produces a list of ToDos from a list of comments
func NewToDos(comments comments.Comments, matchPhrases []string) ToDos {
	todos := make(ToDos, 0)
	for _, comment := range comments {
		todo := NewToDo(*comment, matchPhrases)
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
			// TODO (patrickdevivo) report this error
			continue
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
