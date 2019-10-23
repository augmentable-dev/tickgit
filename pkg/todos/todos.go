package todos

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/augmentable-dev/lege"
	"github.com/src-d/enry/v2"
)

// CStyleCommentOptions ...
var CStyleCommentOptions *lege.ParseOptions = &lege.ParseOptions{
	BoundaryOptions: []lege.BoundaryOption{
		lege.BoundaryOption{
			Starts: []string{"//"},
			Ends:   []string{"\n"},
		},
		lege.BoundaryOption{
			Starts: []string{"/*"},
			Ends:   []string{"*/"},
		},
	},
}

// HashStyleCommentOptions ...
var HashStyleCommentOptions *lege.ParseOptions = &lege.ParseOptions{
	BoundaryOptions: []lege.BoundaryOption{
		lege.BoundaryOption{
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

// SearchFile ...
func SearchFile(filePath string) ([]*lege.Collection, error) {
	src, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	lang := Language(enry.GetLanguage(filepath.Base(filePath), src))
	if enry.IsVendor(filePath) {
		return nil, nil
	}
	options := LanguageParseOptions[lang]
	commentParser, err := lege.NewParser(options)
	if err != nil {
		return nil, err
	}
	comments, err := commentParser.Parse(bytes.NewReader(src))
	if err != nil {
		return nil, err
	}

	for _, comment := range comments {
		s := comment.String()
		if !strings.Contains(s, "TODO") {
			continue
		}
		s = strings.Replace(comment.String(), "TODO", "", 1)
		s = strings.Trim(s, " ")
		fmt.Printf("%q\n", s)
	}

	return comments, nil
}
