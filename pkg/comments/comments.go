package comments

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/augmentable-dev/lege"
	"github.com/go-enry/go-enry/v2"
	"github.com/karrick/godirwalk"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

// Comments is a list of comments
type Comments []*Comment

// Comment represents a comment in a source code file
type Comment struct {
	lege.Collection
	FilePath string
}

// SearchFile searches a file for comments. It infers the language
func SearchFile(filePath string, reader io.Reader, cb func(*Comment)) error {
	// create a preview reader that reads in some of the file for enry to better identify the language
	var buf bytes.Buffer
	tee := io.TeeReader(reader, &buf)
	previewReader := io.LimitReader(tee, 1000)
	preview, err := ioutil.ReadAll(previewReader)
	if err != nil {
		return err
	}

	// create a new reader concatenating the preview and the original reader (which has now been read from)
	fullReader := io.MultiReader(strings.NewReader(buf.String()), reader)

	lang := Language(enry.GetLanguage(filepath.Base(filePath), preview))
	if enry.IsVendor(filePath) {
		return nil
	}
	options, ok := LanguageParseOptions[lang]
	if !ok { // TODO provide a default parse option for when we don't know how to handle a language? I.e. default to CStyle comments say
		return nil
	}
	commentParser, err := lege.NewParser(options)
	if err != nil {
		return err
	}

	collections, err := commentParser.Parse(fullReader)
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
	err = fileIter.ForEach(func(file *object.File) error {
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

	if err != nil {
		return err
	}

	wg.Wait()
	return nil
}
