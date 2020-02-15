## API

### TODOs Badge

`GET` requests to `https://api.tickgit.com/badge` with a `repo` query param:

```
https://api.tickgit.com/badge?repo=https://github.com/facebook/react
```

Supplying a `branch` query param will lookup a specific branch.

Will return JSON that can be fed into a shields.io badge: [https://shields.io/endpoint](https://shields.io/endpoint)

[![TODOs](https://img.shields.io/endpoint?url=https://api.tickgit.com/badge?repo=github.com/augmentable-dev/tickgit)](https://www.tickgit.com/browse?repo=github.com/augmentable-dev/tickgit)
