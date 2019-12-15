package comments

import (
	"testing"
)

func TestJSFiles(t *testing.T) {
	var comments Comments
	err := SearchDir("testdata/javascript", func(comment *Comment) {
		comments = append(comments, comment)
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(comments) != 2 {
		t.Fail()
	}
}

func TestLispFiles(t *testing.T) {
	var comments Comments
	err := SearchDir("testdata/lisp", func(comment *Comment) {
		comments = append(comments, comment)
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(comments) != 1 {
		t.Fail()
	}
}
