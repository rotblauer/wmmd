# Wub

> Delivers live-reloaded markdown files from a directory to your browser,
watching recursively for changes to show you the one you're working on.

![wub-screenshot](./wub.png)

----

The project was begun with editing (a lot of) Github Wiki pages in mind; the style and layout
mimic Github as closely as my patience allows. All rendering is done locally;
you won't need an internet connection, and you'll never have to worry
about hitting a Github API rate limit.

### Installation
```shell
$ go get github.com/rotblauer/wub/...
$ go install github.com/rotblauer/wub
$ which wub
> $GOPATH/bin/wub
```

### Usage
Wub is simple. Point it a directory and run.

It will open to the last edited markdown file in that directory.
Then, just edit any markdown file in that directory or any subdirectory
and wub will detect that change and render it.

Relative links will be functional, i.e. `./Instructions.md` or `./Instructions`.
Relative image resources should work too (wub attempts some cleverness because
Github requires relative image paths to be prefixed with `/wiki`).

__Basic:__
```shell
$ cd my/markdown/directory
$ wub
```

__Options:__
```shell
$ wub [--options] [PATH]
```

| Option | Default | About |
|---|---|---|
| `--port` | 3000 | Specify port to serve on |
| `--topless` | false | Remove leading tags from file |
| `PATH` | `$CWD` | May be an absolute or relative path |

The `--topless` option removes:
```md
---
name: Home
category: Documentation
info: The first instance of these two '---' lines and anything between them will not be rendered.
---
```

#### Wiki mode
Press `w`, or click the light gray button in the top right,
to toggle Github Wiki page style layout, which renders `_Sidebar.*` and
`_Footer.*` in their respective places.

### Limitations and ~~shit~~ hit list
- It does not scroll for you to your current changes... not sure if that's a limitation or a feature.
- You have to `go get` it; wub depends on external html, css, and image files and I don't
 have the patience to transfer them to bindata.
- Otherwise it is perfect.
