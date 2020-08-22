package comments

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestGitIgnore(t *testing.T) {
	gitignorePath := "testdata/gitignore/.gitignore"
	err := ioutil.WriteFile(gitignorePath, []byte("test.go\n"), 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := os.Remove(gitignorePath)
		if err != nil {
			t.Fatal(err)
		}
	}()
	var comments Comments
	err = SearchDir("testdata/gitignore", func(comment *Comment) {
		comments = append(comments, comment)
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(comments) != 0 {
		t.Fail()
	}
}

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

func TestKotlinFiles(t *testing.T) {
	var comments Comments
	err := SearchDir("testdata/kotlin", func(comment *Comment) {
		comments = append(comments, comment)
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(comments) != 2 {
		t.Fail()
	}
}

func TestJuliaFiles(t *testing.T) {
	var comments Comments
	err := SearchDir("testdata/julia", func(comment *Comment) {
		comments = append(comments, comment)
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(comments) != 3 {
		t.Fail()
	}
}

func TestElixirFiles(t *testing.T) {
	var comments Comments
	err := SearchDir("testdata/elixir", func(comment *Comment) {
		comments = append(comments, comment)
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(comments) != 2 {
		t.Fail()
	}
}

func TestHaskellFiles(t *testing.T) {
	var comments Comments
	err := SearchDir("testdata/haskell", func(comment *Comment) {
		comments = append(comments, comment)
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(comments) != 2 {
		t.Fail()
	}
}
