## API

If you'd like to access tickgit information about public `git` repositories, you can use our API.

### TODOs

`GET` requests to `https://tickgit.augmentable.dev/todos` with the `repo` query param populated with the URL of a git repo, like so:

```
https://tickgit.augmentable.dev/todos?repo=https://github.com/facebook/react
```
Will return a simple JSON response:

```
{"todos":125}
```

Indicating the total count of TODOs found in the `HEAD` of that repository.

_more coming soon!_

### TODOs Badge

Similarly, `GET` requests to `https://tickgit.augmentable.dev/todos-badge` with the same `repo` query param:

```
http://tickgit.augmentable.dev/todos-badge?repo=https://github.com/facebook/react
```

Will return JSON that can be fed into a shields.io badge: [https://shields.io/endpoint](https://shields.io/endpoint)
