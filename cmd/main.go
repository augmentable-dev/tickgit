package main

import (
	"os"

	"github.com/augmentable-dev/tickgit"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: "https://github.com/augmentable-dev/tickgit",
	})
	handleError(err)

	ref, err := r.Head()
	handleError(err)

	commit, err := r.CommitObject(ref.Hash())
	handleError(err)

	err = tickgit.WriteStatus(commit, os.Stdout)
	handleError(err)
}
