package comments

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/augmentable-dev/lege"
	"github.com/karrick/godirwalk"
	"github.com/src-d/enry/v2"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

// CStyleCommentOptions ...
var CStyleCommentOptions *lege.ParseOptions = &lege.ParseOptions{
	Boundaries: []lege.Boundary{
		{
			Start: "//",
			End:   "\n",
		},
		{
			Start: "/*",
			End:   "*/",
		},
	},
}

// HashStyleCommentOptions ...
var HashStyleCommentOptions *lege.ParseOptions = &lege.ParseOptions{
	Boundaries: []lege.Boundary{
		{
			Start: "#",
			End:   "\n",
		},
	},
}

// LispStyleCommentOptions ..
var LispStyleCommentOptions *lege.ParseOptions = &lege.ParseOptions{
	Boundaries: []lege.Boundary{
		{
			Start: ";",
			End:   "\n",
		},
	},
}

// Language is a source language (i.e. "Go")
type Language string

// LanguageParseOptions keeps track of source languages and their corresponding comment options
var LanguageParseOptions map[Language]*lege.ParseOptions = map[Language]*lege.ParseOptions{
	"Go":           CStyleCommentOptions,
	"Java":         CStyleCommentOptions,
	"C":            CStyleCommentOptions,
	"C++":          CStyleCommentOptions,
	"C#":           CStyleCommentOptions,
	"JavaScript":   CStyleCommentOptions,
	"Python":       HashStyleCommentOptions,
	"Ruby":         HashStyleCommentOptions,
	"PHP":          CStyleCommentOptions,
	"Shell":        HashStyleCommentOptions,
	"Visual Basic": {Boundaries: []lege.Boundary{{Start: "'", End: "\n"}}},
	"TypeScript":   CStyleCommentOptions,
	"Objective-C":  CStyleCommentOptions,
	"Groovy":       CStyleCommentOptions,
	"Swift":        CStyleCommentOptions,
	"Common Lisp":  LispStyleCommentOptions,
	"Emacs Lisp":   LispStyleCommentOptions,
	"R":            HashStyleCommentOptions,
}

// Comments is a list of comments
type Comments []*Comment

// Comment represents a comment in a source code file
type Comment struct {
	lege.Collection
	FilePath string
}

// SearchFile searches a file for comments. It infers the language
func SearchFile(filePath string, reader io.Reader, cb func(*Comment)) error {
	// TODO right now, enry only infers the language based on the file extension
	// we should add some "preview" bytes from the file so that it has some sample content to examine
	lang := Language(enry.GetLanguage(filepath.Base(filePath), nil))
	if enry.IsVendor(filePath) {
		return nil
	}
	options, ok := LanguageParseOptions[lang]
	if !ok { // TODO provide a default parse option?
		return nil
	}
	commentParser, err := lege.NewParser(options)
	if err != nil {
		return err
	}

	collections, err := commentParser.Parse(reader)
	if err != nil {
		return err
	}

	for _, c := range collections {
		comment := Comment{*c, filePath}
		cb(&comment)
	}

	return nil
}

// SearchDir searches a directory for comments
func SearchDir(dirPath string, cb func(comment *Comment)) error {
	err := godirwalk.Walk(dirPath, &godirwalk.Options{
		Callback: func(path string, de *godirwalk.Dirent) error {
			localPath, err := filepath.Rel(dirPath, path)
			if err != nil {
				return err
			}
			pathComponents := strings.Split(localPath, string(os.PathSeparator))
			// let's ignore git directories TODO: figure out a more generic way to set ignores
			matched, err := filepath.Match(".git", pathComponents[0])
			if err != nil {
				return err
			}
			if matched {
				return nil
			}
			if de.IsRegular() {
				p, err := filepath.Abs(path)
				if err != nil {
					return err
				}
				f, err := os.Open(p)
				if err != nil {
					return err
				}
				err = SearchFile(localPath, f, cb)
				if err != nil {
					return err
				}
				f.Close()
			}
			return nil
		},
		Unsorted: true,
	})
	if err != nil {
		return err
	}
	return nil
}

// SearchCommit searches all files in the tree of a given commit
func SearchCommit(commit *object.Commit, cb func(*Comment)) error {
	var wg sync.WaitGroup
	errs := make(chan error)

	fileIter, err := commit.Files()
	if err != nil {
		return err
	}
	defer fileIter.Close()
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
				err = SearchFile(file.Name, r, cb)
				if err != nil {
					errs <- err
					return
				}

			}()
		}
		return nil
	})

	wg.Wait()
	return nil
}
