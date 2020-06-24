package todos

import (
	"testing"

	"github.com/augmentable-dev/lege"
	"github.com/augmentable-dev/tickgit/pkg/comments"
)

func TestNewToDoNil(t *testing.T) {
	collection := lege.NewCollection(lege.Location{}, lege.Location{}, lege.Boundary{}, "Hello World")
	comment := comments.Comment{
		Collection: *collection,
	}
	todo := NewToDo(comment, startingMatchPhrases)

	if todo != nil {
		t.Fatalf("did not expect a TODO, got: %v", todo)
	}
}

func TestNewToDo(t *testing.T) {
	collection := lege.NewCollection(lege.Location{}, lege.Location{}, lege.Boundary{}, "TODO Hello World")
	comment := comments.Comment{
		Collection: *collection,
	}
	todo := NewToDo(comment, startingMatchPhrases)

	if todo == nil {
		t.Fatalf("expected a TODO, got: %v", todo)
	}

	if todo.Phrase != "TODO" {
		t.Fatalf("expected matched phrase to be TODO, got: %s", todo.Phrase)
	}
}
