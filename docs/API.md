## API

### TODOs Badge

`GET` requests to `https://api.tickgit.com/badgen` with a `repo` path segment:

```
https://api.tickgit.com/badgen/github.com/facebook/react
```

Supplying a `branch` segment will lookup a specific branch. `master` is the branch used if none is specified.

```
https://api.tickgit.com/badgen/github.com/facebook/react/branch-name
```

Will return JSON that can be fed into a badgen badge: [https://badgen.net/https](https://badgen.net/https)

```
[![TODOs](https://badgen.net/https/api.tickgit.com/badgen/github.com/augmentable-dev/tickgit)](https://www.tickgit.com/browse?repo=github.com/augmentable-dev/tickgit)
```

[![TODOs](https://badgen.net/https/api.tickgit.com/badgen/github.com/augmentable-dev/tickgit)](https://www.tickgit.com/browse?repo=github.com/augmentable-dev/tickgit)
