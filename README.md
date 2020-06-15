[![GoDoc](https://godoc.org/github.com/augmentable-dev/tickgit?status.svg)](https://godoc.org/github.com/augmentable-dev/tickgit)
[![BuildStatus](https://github.com/augmentable-dev/tickgit/workflows/tests/badge.svg)](https://github.com/augmentable-dev/tickgit/actions?workflow=tests)
[![Go Report Card](https://goreportcard.com/badge/github.com/augmentable-dev/tickgit)](https://goreportcard.com/report/github.com/augmentable-dev/tickgit)
![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/augmentable-dev/tickgit)
[![Coverage](http://gocover.io/_badge/github.com/augmentable-dev/tickgit)](http://gocover.io/github.com/augmentable-dev/tickgit)
[![TODOs](https://badgen.net/https/api.tickgit.com/badgen/github.com/augmentable-dev/tickgit)](https://www.tickgit.com/browse?repo=github.com/augmentable-dev/tickgit)

## tickgit 🎟️

`tickgit` is a tool to help you manage latent work in a codebase. Use the `tickgit` command to view pending tasks, progress reports, completion summaries and historical data (using `git` history).

It's not meant to replace full-fledged project management tools such as JIRA or Trello. It will, hopefully, be a useful way to augment those tools with project management patterns that coexist with your code. As such, it's primary audience is software engineers.

### TODOs

`tickgit` will scan a codebase and identify any TODO items in the comments. It will output a report like so:

```
# tickgit ~/Desktop/facebook/react
...
TODO:
  => packages/scheduler/src/__tests__/SchedulerBrowser-test.js:85:9
  => added 1 month ago by Andrew Clark <git@andrewclark.io> in a2e05b6c148b25590884e8911d4d4acfcb76a487

TODO: Scheduler no longer requires these methods to be polyfilled. But
  => packages/scheduler/src/__tests__/SchedulerBrowser-test.js:77:7
  => added 1 month ago by Andrew Clark <git@andrewclark.io> in a2e05b6c148b25590884e8911d4d4acfcb76a487

TODO: Scheduler no longer requires these methods to be polyfilled. But
  => packages/scheduler/src/forks/SchedulerHostConfig.default.js:77:7
  => added 1 month ago by Andrew Clark <git@andrewclark.io> in a2e05b6c148b25590884e8911d4d4acfcb76a487

TODO: useTransition hook instead.
  => fixtures/concurrent/time-slicing/src/index.js:110:11
  => added 3 weeks ago by Sebastian Markbåge <sebastian@calyptus.eu> in 3ad076472ce9108b9b8a6a6fe039244b74a34392

128 TODOs Found 📝
```

Check out [an example](https://www.tickgit.com/browse?repo=github.com/kubernetes/kubernetes) of the TODOs tickgit will surface for the Kubernetes codebase.

#### Coming Soon

- [x] Blame - get a better sense of how old TODOs are, when they were introduced and by whom
- [x] More `TODO` type phrases to match, such as `FIXME`, `XXX`, `HACK`, or customized alternatives.
- [ ] Context - more visibility into the lines of code _around_ a TODO for greater context
- [ ] More configurability (e.g. custom ignore paths)
- [ ] Markdown parsing
- [ ] More thorough historical stats

### Installation

#### Homebrew

```
brew tap augmentable-dev/tickgit
brew install tickgit
```

#### GoBinaries

```
curl -sf https://gobinaries.com/augmentable-dev/tickgit/cmd/tickgit | sh
```

Will use [GoBinaries](https://gobinaries.com/) to install the latest version of `tickgit`.
You can specifiy a particular version by appending `@VERSION` to the URL above.

#### go install

```
go install github.com/augmentable-dev/tickgit/cmd/tickgit
```

With `$GOBIN` set and in your `$PATH`.

### Usage

The most up to date usage will be the output of `tickgit --help`.

### API

To find information about using the tickgit API, see [this file](https://github.com/augmentable-dev/tickgit/blob/master/docs/API.md).
