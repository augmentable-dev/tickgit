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

	if len(comments) != 3 {
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

func TestRustFiles(t *testing.T) {
	var comments Comments
	err := SearchDir("testdata/rust", func(comment *Comment) {
		comments = append(comments, comment)
	})
	if err != nil {
		t.Fatal(err)
	}

	// TODO: break the different comment types out into separate files?
	// once the issue with lege is worked out for handling the different comment types
	if len(comments) != 21 {
		t.Fail()
	}
}

func TestPHPFiles(t *testing.T) {
	var comments Comments
	err := SearchDir("testdata/php", func(comment *Comment) {
		comments = append(comments, comment)
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(comments) != 3 {
		t.Fail()
	}
}
