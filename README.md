[![GoDoc](https://godoc.org/github.com/augmentable-dev/tickgit?status.svg)](https://godoc.org/github.com/augmentable-dev/tickgit)
[![BuildStatus](https://github.com/augmentable-dev/tickgit/workflows/tests/badge.svg)](https://github.com/augmentable-dev/tickgit/actions?workflow=tests)
[![Go Report Card](https://goreportcard.com/badge/github.com/augmentable-dev/tickgit)](https://goreportcard.com/report/github.com/augmentable-dev/tickgit)
![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/augmentable-dev/tickgit)
[![Coverage](http://gocover.io/_badge/github.com/augmentable-dev/tickgit)](http://gocover.io/github.com/augmentable-dev/tickgit)
[![TODOs](https://img.shields.io/endpoint?url=https%3A%2F%2Ftickgit.augmentable.dev%2Ftodos-badge%3Frepo%3Dhttps%3A%2F%2Fgithub.com%2Faugmentable-dev%2Ftickgit)](https://todos.augmentable.dev/?repo=https://github.com/augmentable-dev/tickgit)

## tickgit 🎟️

`tickgit` is a tool to help you manage tickets, todo items, and checklists within a codebase. Use the `tickgit` command to view pending tasks, progress reports, completion summaries and historical data (using `git` history).

It's not meant to replace full-fledged project management tools such as JIRA or Trello. It will, hopefully, be a useful way to augment those tools with project management patterns that coexist with your code. As such, it's primary audience is software engineers.

### TODOs

`tickgit todos` will scan a codebase and identify any TODO items in the comments. It will output a report like so:

```
# tickgit todos ~/Desktop/facebook/react
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

#### Coming Soon

- [x] History - get a better sense of how old TODOs are, when they were introduced and by whom
- [ ] Context - more visibility into the lines of code _around_ a TODO for greater context

### Tickets

Tickets are a way of defining more complex tasks in your codebase as config files. Currently, tickets are HCL files that look like the following:

```hcl
# rocketship.tickgit

goal "Build the Rocketship 🚀" {
    description = "Finalize the construction of the Moonblaster 2000"

    task "Construct the engines" {
        status = "done"
    }

    task "Attach the engines" {
        status = "pending"
    }

    task "Thoroughly test the engines" {
        status = "pending"
    }
}
```

```
$ tickgit status
=== Build the Rocketship 🚀 ⏳
  --- 1/3 tasks completed (2 remaining)
  --- 33% completed

  ✅ Construct the engines
  ⏳ Attach the engines
  ⏳ Thoroughly test the engines
```

#### Coming Soon

- [ ] Simpler ticket definitions - in YAML and/or other (less verbose) config languages
- [ ] More complex tickets - more states, dependencies on other tickets, etc

### Checklists

_Coming soon_. Checklists will be a way of parsing Markdown checklists in your codebase (either in `.md` files, or within your comments).


### Why is this useful?

This project is a proof-of-concept. Keeping tickets next to the code they're meant to describe could have the following benefits:

- Tickets live with the code, no need for a 3rd party tool or system (anyone with git access to the repository has access to contributing to the tickets)
- Updating a ticket's status and merging/committing code are the same action, no need to synchronize across multiple tools
- Source of truth for a project's ticket history is now the git history, which can be queried and analyzed
- Current status of a `goal` can be reported by simply parsing the repository's `head`
- Less context switching between the codebase itself and the system describing "what needs to be done"

Generally speaking, this is an experiment in ways to do project management, within the codebase of a project. With a `git` history and some clever parsing, quite a bit of metadata about a project can be gleaned from its codebase. Let's see how useful we can make that information.

### Installation

#### Homebrew

```
brew tap augmentable-dev/tickgit
brew install tickgit
```


### Usage

The most up to date usage will be the output of `tickgit --help`. The most common usage, however, is `tickgit status` which will print a status report of tickets for a given git repository. By default, it uses the current working directory.

### API

To find information about using the tickgit API, see [this file](https://github.com/augmentable-dev/tickgit/blob/master/docs/API.md).
