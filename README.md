## tickgit

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

  ✅ Construct the engines:
  ⏳ Attach the engines:
  ⏳ Thoroughly test the engines:
```
