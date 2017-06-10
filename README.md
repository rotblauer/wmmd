# Wub

> Delivers live-reloaded markdown files from a directory to your browser,
watching recursively for changes to show you the one you're working on.

![wub-screenshot](./wub.png)

The project was begun with editing Github Wiki pages in mind. CSS style and layout
mimic Github as closely as my patience allows.

### Installation
```shell
$ go get github.com/rotblauer/wub/...
$ go install github.com/rotblauer/wub
$ which wub
> $GOPATH/bin/wub
```

### Usage
Wub is simple. Point it a directory and run.

__Basic:__
```shell
$ cd my/markdown/directory
$ wub
```

__Advanced:__
```shell
$ wub --port 3001 my/markdown/directory
```

It will open to the last edited markdown file.
Then, just edit any markdown file in that directory or any subdirectory and wub will detect that change and render it.

Relative links will be functional, i.e. `./Instructions.md`. Image resources relative to the base directory will be functional.

#### Wiki mode
Press `w`, or click the light gray button "Wiki: [on|off]" in the top right
to toggle Github Wiki page style layout, which renders `_Sidebar.md` and
`_Footer.md` in their respective places.

### Limitations and ~~shit~~ hit list
- It is not very clever about file names and titles (with regard to links and titles); it doesn't handle ambiguity well.
As far as it goes so far is appending ".md" to links without that extension, e.g. it would be nice if, upon receiving a
url parameter `./Instructions`, it would look for `.md`, `.markdown`, `.mdown`, `.adoc`, instead of just `.md`
- It would be really great if it were clever about absolute href/urls, so "hardcoded" wiki links could
be toggled to relative paths if they exits. Not _yet_.
- It does not scroll for you to your current changes... not sure if that's a limitation or a feature.
