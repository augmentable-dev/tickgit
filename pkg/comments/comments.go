package comments

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/augmentable-dev/lege"
	"github.com/src-d/enry/v2"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

// CStyleCommentOptions ...
var CStyleCommentOptions *lege.ParseOptions = &lege.ParseOptions{
	Boundaries: []lege.Boundary{
		lege.Boundary{
			Starts: []string{"//"},
			Ends:   []string{"\n"},
		},
		lege.Boundary{
			Starts: []string{"/*"},
			Ends:   []string{"*/"},
		},
	},
}

// HashStyleCommentOptions ...
var HashStyleCommentOptions *lege.ParseOptions = &lege.ParseOptions{
	Boundaries: []lege.Boundary{
		lege.Boundary{
			Starts: []string{"#"},
			Ends:   []string{"\n"},
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

// SearchDir searches a directory for comments
func SearchDir(dirPath string) (Comments, error) {
	found := make(Comments, 0)
	// TODO let's see what we can do concurrently here to speed up the processing
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
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
		if !info.IsDir() {
			p, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			f, err := os.Open(p)
			if err != nil {
				return err
			}
			t, err := SearchFile(p, f)
			if err != nil {
				return err
			}
			c := Comments(t)
			found = append(found, c...)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return found, nil
}

// SearchCommit searches all files in the tree of a given commit
func SearchCommit(commit *object.Commit) (Comments, error) {
	c := make(Comments, 0)
	t, err := commit.Tree()
	if err != nil {
		return nil, err
	}
	fileIter := t.Files()
	fileIter.ForEach(func(file *object.File) error {
		if file.Mode.IsFile() {
			r, err := file.Reader()
			if err != nil {
				return err
			}
			found, err := SearchFile(file.Name, r)
			if err != nil {
				return err
			}
			c = append(c, found...)
		}
		return nil
	})
	return c, nil
}
