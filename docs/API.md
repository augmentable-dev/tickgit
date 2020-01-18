## API

### TODOs Badge

`GET` requests to `https://api.tickgit.com/todos-badge` with a `repo` query param:

```
http://api.tickgit.com/todos-badge?repo=https://github.com/facebook/react
```

Supplying a `branch` query param will lookup a specific branch.

Will return JSON that can be fed into a shields.io badge: [https://shields.io/endpoint](https://shields.io/endpoint)
