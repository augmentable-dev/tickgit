package comments

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/augmentable-dev/lege"
	"github.com/src-d/enry/v2"
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
func SearchFile(filePath string) (Comments, error) {
	src, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	lang := Language(enry.GetLanguage(filepath.Base(filePath), src))
	if enry.IsVendor(filePath) {
		return nil, nil
	}
	options, ok := LanguageParseOptions[lang]
	if !ok {
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
			t, err := SearchFile(p)
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
