package comments

import (
	"bytes"
	"io"
	"io/ioutil"
	"path/filepath"
	"sync"

	"github.com/augmentable-dev/lege"
	"github.com/src-d/enry/v2"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

// CStyleCommentOptions ...
var CStyleCommentOptions *lege.ParseOptions = &lege.ParseOptions{
	Boundaries: []lege.Boundary{
		lege.Boundary{
			Start: "//",
			End:   "\n",
		},
		lege.Boundary{
			Start: "/*",
			End:   "*/",
		},
	},
}

// HashStyleCommentOptions ...
var HashStyleCommentOptions *lege.ParseOptions = &lege.ParseOptions{
	Boundaries: []lege.Boundary{
		lege.Boundary{
			Start: "#",
			End:   "\n",
		},
	},
}

// Language is a source language (i.e. "Go")
type Language string

// LanguageParseOptions keeps track of source languages and their corresponding comment options
var LanguageParseOptions map[Language]*lege.ParseOptions = map[Language]*lege.ParseOptions{
	"Go":         CStyleCommentOptions,
	"Java":       CStyleCommentOptions,
	"C":          CStyleCommentOptions,
	"C++":        CStyleCommentOptions,
	"C#":         CStyleCommentOptions,
	"JavaScript": CStyleCommentOptions,
	"Python":     HashStyleCommentOptions,
	"Ruby":       HashStyleCommentOptions,
	"PHP":        CStyleCommentOptions,
}

// Comments is a list of comments
type Comments []Comment

// Comment represents a comment in a source code file
type Comment struct {
	lege.Collection
	FilePath string
}

// SearchFile searches a file for comments. It infers the language
func SearchFile(filePath string, reader io.ReadCloser) (Comments, error) {
	src, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	lang := Language(enry.GetLanguage(filepath.Base(filePath), src))
	if enry.IsVendor(filePath) {
		return nil, nil
	}
	options, ok := LanguageParseOptions[lang]
	if !ok { // TODO provide a default parse option?
		return nil, nil
	}
	commentParser, err := lege.NewParser(options)
	if err != nil {
		return nil, err
	}

	collections, err := commentParser.Parse(bytes.NewReader(src))
	if err != nil {
		return nil, err
	}

	comments := make(Comments, 0)
	for _, c := range collections {
		comment := Comment{*c, filePath}
		comments = append(comments, comment)
	}

	return comments, nil
}

// SearchCommit searches all files in the tree of a given commit
func SearchCommit(commit *object.Commit) (Comments, error) {
	found := make(Comments, 0)
	t, err := commit.Tree()
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	errs := make(chan error)

	fileIter := t.Files()
	fileIter.ForEach(func(file *object.File) error {
		if file.Mode.IsFile() {
			wg.Add(1)
			go func() {
				defer wg.Done()

				r, err := file.Reader()
				if err != nil {
					errs <- err
					return
				}
				c, err := SearchFile(file.Name, r)
				if err != nil {
					errs <- err
					return
				}

				for _, comment := range c {
					found = append(found, comment)
				}
			}()
		}
		return nil
	})

	wg.Wait()
	return found, nil
}
