package tickgit

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

func loadFixtureRepo(repoName string) (*git.Repository, error) {
	r, err := git.Init(memory.NewStorage(), memfs.New())
	if err != nil {
		return nil, err
	}
	w, err := r.Worktree()
	if err != nil {
		return nil, err
	}

	fixturesRepoDir := filepath.Join("testdata/repos", repoName)
	commits, err := ioutil.ReadDir(fixturesRepoDir)
	if err != nil {
		return nil, err
	}

	walkFunc := func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			f, err := w.Filesystem.Create(path)
			if err != nil {
				return err
			}

			contents, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			_, err = f.Write(contents)
			if err != nil {
				return err
			}
			w.Add(path)
		}
		return nil
	}

	for _, commit := range commits {
		// add all the files
		err := filepath.Walk(filepath.Join(fixturesRepoDir, commit.Name()), walkFunc)
		if err != nil {
			return nil, err
		}
		_, err = w.Commit(commit.Name(), &git.CommitOptions{
			Author: &object.Signature{
				Name:  "User1",
				Email: "user1@example.com",
				When:  time.Now(),
			},
		})
		if err != nil {
			return nil, err
		}
	}

	return r, nil
}

func TestParse(t *testing.T) {
	file, err := Parse([]byte(`
	goal "My Long Term Goal" {
		task "Step 1" {
			status = "pending"
		}
	}
	`), "test.tickgit")

	if err != nil {
		t.Fatal(err)
	}

	goals := file.Goals
	{
		want := 1
		got := len(goals)
		if want != got {
			t.Fatalf("unexpected number of goals, want: %d got: %d", want, got)
		}
	}
}

func TestGit(t *testing.T) {
	r, err := loadFixtureRepo("repo-001")
	if err != nil {
		t.Fatal(err)
	}

	latest, err := r.Head()
	if err != nil {
		t.Fatal(err)
	}

	commit, err := r.CommitObject(latest.Hash())
	if err != nil {
		t.Fatal(err)
	}

	goals, err := GoalsFromCommit(commit, nil)
	if err != nil {
		t.Fatal(err)
	}

	{
		want := 2
		got := len(goals)
		if want != got {
			t.Fatalf("unexpected number of goals, want: %d got: %d", want, got)
		}
	}

	{
		want := 3
		got := len(goals[0].Tasks)
		if want != got {
			t.Fatalf("unexpected number of tasks in first goal, want: %d got: %d", want, got)
		}
	}

}

func TestPrintStatus(t *testing.T) {
	r, err := loadFixtureRepo("repo-001")
	if err != nil {
		t.Fatal(err)
	}

	latest, err := r.Head()
	if err != nil {
		t.Fatal(err)
	}

	commit, err := r.CommitObject(latest.Hash())
	if err != nil {
		t.Fatal(err)
	}

	err = WriteStatus(commit, os.Stdout)
	if err != nil {
		t.Fatal(err)
	}
}
