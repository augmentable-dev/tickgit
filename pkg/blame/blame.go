package blame

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Options are options to determine what and how to blame
type Options struct {
	Directory string
	SHA       string
	Lines     []int
}

// Blame represents the "blame" of a particlar line or range of lines
type Blame struct {
	SHA       string
	Author    Event
	Committer Event
	Range     [2]int
}

// Event represents the who and when of a commit event
type Event struct {
	Name  string
	Email string
	When  time.Time
}

func (blame *Blame) String() string {
	return fmt.Sprintf("%s: %s <%s>", blame.SHA, blame.Author.Name, blame.Author.Email)
}

func (event *Event) String() string {
	return fmt.Sprintf("%s <%s>", event.Name, event.Email)
}

// Result is a mapping of line numbers to blames for a given file
type Result map[int]Blame

func (options *Options) argsFromOptions(filePath string) []string {
	args := []string{"blame"}
	if options.SHA != "" {
		args = append(args, options.SHA)
	}

	for _, line := range options.Lines {
		args = append(args, fmt.Sprintf("-L %d,%d", line, line))
	}

	args = append(args, "--porcelain", "--incremental")

	args = append(args, filePath)
	return args
}

func parsePorcelain(reader io.Reader) (Result, error) {
	scanner := bufio.NewScanner(reader)
	res := make(Result)

	const (
		author     = "author "
		authorMail = "author-mail "
		authorTime = "author-time "
		authorTZ   = "author-tz "

		committer     = "committer "
		committerMail = "committer-mail "
		committerTime = "committer-time "
		committerTZ   = "committer-tz "
	)

	seenCommits := make(map[string]Blame)
	var currentCommit Blame
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case strings.HasPrefix(line, author):
			currentCommit.Author.Name = strings.TrimPrefix(line, author)
		case strings.HasPrefix(line, authorMail):
			s := strings.TrimPrefix(line, authorMail)
			currentCommit.Author.Email = strings.Trim(s, "<>")
		case strings.HasPrefix(line, authorTime):
			timeString := strings.TrimPrefix(line, authorTime)
			i, err := strconv.ParseInt(timeString, 10, 64)
			if err != nil {
				return nil, err
			}
			currentCommit.Author.When = time.Unix(i, 0)
		case strings.HasPrefix(line, authorTZ):
			tzString := strings.TrimPrefix(line, authorTZ)
			parsed, err := time.Parse("-0700", tzString)
			if err != nil {
				return nil, err
			}
			loc := parsed.Location()
			currentCommit.Author.When = currentCommit.Author.When.In(loc)
		case strings.HasPrefix(line, committer):
			currentCommit.Committer.Name = strings.TrimPrefix(line, committer)
		case strings.HasPrefix(line, committerMail):
			s := strings.TrimPrefix(line, committer)
			currentCommit.Committer.Email = strings.Trim(s, "<>")
		case strings.HasPrefix(line, committerTime):
			timeString := strings.TrimPrefix(line, committerTime)
			i, err := strconv.ParseInt(timeString, 10, 64)
			if err != nil {
				return nil, err
			}
			currentCommit.Committer.When = time.Unix(i, 0)
		case strings.HasPrefix(line, committerTZ):
			tzString := strings.TrimPrefix(line, committerTZ)
			parsed, err := time.Parse("-0700", tzString)
			if err != nil {
				return nil, err
			}
			loc := parsed.Location()
			currentCommit.Committer.When = currentCommit.Committer.When.In(loc)
		case len(strings.Split(line, " ")[0]) == 40: // if the first string sep by a space is 40 chars long, it's probably the commit header
			split := strings.Split(line, " ")
			sha := split[0]

			// if we haven't seen this commit before, create an entry in the seen commits map that will get filled out in subsequent lines
			if _, ok := seenCommits[sha]; !ok {
				seenCommits[sha] = Blame{SHA: sha}
			}

			// update the current commit to be this new one we've just encountered
			currentCommit.SHA = sha

			// pull out the line information
			line := split[2]
			l, err := strconv.ParseInt(line, 10, 64) // the starting line of the range
			if err != nil {
				return nil, err
			}

			var c int64
			if len(split) > 3 {
				c, err = strconv.ParseInt(split[3], 10, 64) // the number of lines in the range
				if err != nil {
					return nil, err
				}
			}
			for i := l; i < l+c; i++ {
				res[int(i)] = Blame{SHA: sha}
			}
		}
		// after every line, make sure the current commit in the seen commits map is updated
		seenCommits[currentCommit.SHA] = currentCommit
	}
	for line, blame := range res {
		res[line] = seenCommits[blame.SHA]
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

// Exec uses git to lookup the blame of a file, given the supplied options
func Exec(ctx context.Context, filePath string, options *Options) (Result, error) {
	gitPath, err := exec.LookPath("git")
	if err != nil {
		return nil, fmt.Errorf("could not find git: %w", err)
	}

	args := options.argsFromOptions(filePath)

	cmd := exec.CommandContext(ctx, gitPath, args...)
	cmd.Dir = options.Directory

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	res, err := parsePorcelain(stdout)
	if err != nil {
		return nil, err
	}

	if err := cmd.Wait(); err != nil {
		return nil, err
	}

	return res, nil
}
