[![GoDoc](https://godoc.org/github.com/augmentable-dev/tickgit?status.svg)](https://godoc.org/github.com/augmentable-dev/tickgit)
[![BuildStatus](https://github.com/augmentable-dev/tickgit/workflows/Test/badge.svg)](https://github.com/augmentable-dev/tickgit/actions?workflow=Test)

## tickgit 🎟️

Tickets as config. Manage your repository's tickets and todo items as configuration files in your codebase. Use the `tickgit` command to view a snapshot of pending items, summaries of completed items, and historical reports of progress using `git` history.


```hcl
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
=== Build the Rocketship 🚀 ⏳
  --- 1/3 tasks completed (2 remaining)
  --- 33% completed

  ✅ Construct the engines
  ⏳ Attach the engines
  ⏳ Thoroughly test the engines
```

### Why is this useful?

To be honest, I'm not sure yet. This project is a POC I'm exploring. Keeping tickets next to the code they're meant to describe could have the following benefits:

- Tickets live with the code, no need for a 3rd party tool or system (anyone with git access to the repository has access to contributing to the tickets)
- Updating a ticket's status and merging/committing code are the same action, no need to synchronize across multiple tools
- Source of truth for a project's ticket history is now the git history, which can be queried and analyzed
- Current status of a `goal` can be reported by simply parsing the repository's `head`
- Less context switching between the codebase itself and the system describing "what needs to be done"