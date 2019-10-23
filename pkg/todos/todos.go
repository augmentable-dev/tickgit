package todos

import (
	"bytes"
	"io/ioutil"
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

// ToDo represents a ToDo item
type ToDo struct {
	FilePath string
	Line     int
	Position int
	String   string
}

// SearchFile searches a file for comments. It infers the language
func SearchFile(filePath string) ([]*ToDo, error) {
	todos := make([]*ToDo, 0)
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
		// fmt.Printf("%q\n", s)
		todo := &ToDo{
			FilePath: filePath,
			Line:     comment.StartLocation.Line,
			Position: comment.StartLocation.Pos,
			String:   s,
		}
		todos = append(todos, todo)
	}

	return todos, nil
}
