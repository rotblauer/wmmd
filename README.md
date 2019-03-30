# WMMD: Weapon of Mass Markdown


> Delivers live-reloaded markdown files from a directory to your browser,
watching recursively for changes to render the one you're working on.

![wmd-gif](./wmd.gif)

----

The project was begun with editing (a lot of) Github Wiki pages in mind; the style and layout
mimic Github as closely as my patience allows. All rendering is done locally;
you won't need an internet connection, and you'll never have to worry
about hitting a Github API rate limit.

### Installation
```shell
$ go get github.com/rotblauer/wmmd/...
$ go install github.com/rotblauer/wmmd
$ which wmmd
> $GOPATH/bin/wmmd
```

### Usage
Wmmd is simple. Point it a directory and run.

It will open to the last edited markdown file in that directory.
Then, just edit any markdown file in that directory or any subdirectory
and wmmd will detect that change and render it.

Relative links will be functional, i.e. `./Instructions.md` or `./Instructions`.
Relative image resources should work too (wmmd attempts some cleverness because
Github requires relative image paths to be prefixed with `/wiki`).

__Basic:__
```shell
$ cd my/markdown/directory
$ wmmd
```

__Options:__
```shell
$ wmmd [--options] [PATH]
```

| Option | Default | About |
|---|---|---|
| `--port` | 3000 | Specify port to serve on |
| `--topless` | false | Remove leading tags from file |
| `-n` | false | Enable hard line breaks for markdown conversion |
| `-s` | true | Scroll to latest changes automatically |
| `PATH` | `$CWD` | May be an absolute or relative path |

The `--topless` option removes anything like:
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

#### Text mode
Press `t`, or click the light gray button in the top right, to toggle typwriter-like styling, replacing the standard Github flavored appearance.

### Limitations and ~~shit~~ hit list
- Sometimes the auto-scroll gets weird.

- You have to `go get` it; wmmd depends on external html, css, and image files and I don't
 have the patience to transfer them to bindata.
