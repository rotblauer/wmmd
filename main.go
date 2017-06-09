package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/fsnotify/fsnotify"
	"github.com/labstack/echo"
	"github.com/olahol/melody"
	"github.com/shurcooL/github_flavored_markdown"
	"strings"
)

var port int
var dirPath string
var currentFile string

type FileContent struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

func setCurrentFile(path string) {
	currentFile = path
}

func getCurrentFile() string {
	return currentFile
}

func init() {
	flag.IntVar(&port, "port", 3000, "port to serve on")
}

func main() {
	flag.Parse()
	dirPath = mustMakeDirPath()
	mm := melody.New()

	currentFile = filepath.Join(dirPath, "README.md")
	if _, e := os.Stat(currentFile); e != nil && os.IsNotExist(e) {
		currentFile = filepath.Join(dirPath, "Home.md")
		if _, ee := os.Stat(currentFile); ee != nil && os.IsNotExist(ee) {
			fs, _ := ioutil.ReadDir(dirPath)
			for _, f := range fs {
				ext := filepath.Ext(f.Name())
				if !strings.HasPrefix(f.Name(), "_") && (ext == ".md" || ext == ".markdown" || ext == ".mdown") {
					currentFile = f.Name()
					break
				}
			}
		}
	}

	mm.HandleConnect(func(s *melody.Session) {
		log.Println("session connected")
		curFile, e := getReadFile(getCurrentFile())
		if e != nil {
			log.Println(e)
		}
		sidebar, e := getReadFile(filepath.Join(dirPath, "_Sidebar.md"))
		if e != nil {
			log.Println(e)
		}
		footer, e := getReadFile(filepath.Join(dirPath, "_Footer.md"))
		if e != nil {
			log.Println(e)
		}
		for _, f := range []FileContent{sidebar, footer, curFile} {
			if (f == FileContent{}) {
				continue
			}
			j, e := json.Marshal(f)
			if e != nil {
				log.Println(e)
				continue
			}
			mm.Broadcast(j)
		}

		j, e := json.Marshal(curFile)
		if e != nil {
			log.Println(e)
			return
		}
		mm.Broadcast(j)
	})
	mm.HandleDisconnect(func(s *melody.Session) {
		log.Println("session disconnected")
	})

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Println("event:", event)
				// if event.Op&fsnotify.Write == fsnotify.Write {
				log.Println("modified file:", event.Name)
				if ff := filepath.Ext(event.Name); ff != "" && ff == ".md" && ff == ".markdown" && ff == ".mdown" {
					log.Println("not markdown file, continuing...")
					continue
				}
				f := getFilePathFromParam(event.Name)
				setCurrentFile(f)
				m, e := getReadFile(f)
				if e != nil {
					log.Println(e)
					continue
				}
				b, e := json.Marshal(m)
				if e != nil {
					log.Println(e)
					continue
				}
				log.Println("broadcasting", string(b), mm)
				mm.Broadcast(b)
				// }
			case err := <-watcher.Errors:
				log.Println("watcher error:", err)
			}
		}
	}()
	// Sometimes on refresh watcher stops silently.
	err = watcher.Add(dirPath)
	if err != nil {
		log.Fatal(err)
	}

	// Echo is polite because it prioritizes these paths, so they can be overlapping,
	// ie. ":filename" overlaps everything except /
	r := echo.New()
	r.File("/", filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "rotblauer", "wub", "index.html"))
	// Static assets.
	r.Static("/assets", filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "rotblauer", "wub", "assets"))
	r.Static("/node_modules/primer-css/build", filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "rotblauer", "wub", "node_modules/primer-css/build"))
	// Websocket.
	r.GET("/x/0", func(c echo.Context) error {
		mm.HandleRequest(c.Response(), c.Request())
		return nil
	})
	// Any other filename.
	r.Any("/:filename", func(c echo.Context) error {
		filename := getFilePathFromParam(c.Param("filename"))
		setCurrentFile(filename)
		c.Response().Header().Set("Cache-Control: no-cache", "true")
		return c.Redirect(http.StatusMovedPermanently, "/")
	})

	log.Println("Listening...", port)
	r.Logger.Fatal(r.Start(":" + strconv.Itoa(port)))
}

func mustMakeDirPath() string {
	args := flag.Args()
	if len(args) == 0 {
		p, e := os.Getwd()
		if e != nil {
			panic(e)
		}
		return p
	}
	abs, e := filepath.Abs(args[0])
	if e != nil {
		panic(e)
	}
	di, de := os.Stat(abs)
	if de != nil {
		panic(de)
	}
	if !di.IsDir() {
		panic("path must be a dir")
	}
	return abs
}

func getFilePathFromParam(param string) string {
	filename := param
	if filename == "" {
		return ""
	}
	// TODO: handle ambiguous filenames in url param
	if !(filepath.Ext(filename) == ".md") {
		filename = filename + ".md"
	}

	if !filepath.IsAbs(filename) {
		filename, _ = filepath.Abs(filename)
	}
	return filename
}

func getReadFile(path string) (FileContent, error) {
	fileBytes, e := ioutil.ReadFile(path)
	if e != nil {
		log.Println(e)
		return FileContent{}, e
	}
	return FileContent{
		Title: filepath.Base(path), // TODO parse File-Name.md syntax => File Name
		// Body:  emoji.Emojitize(string(github_flavored_markdown.Markdown(fileBytes))),
		Body: string(github_flavored_markdown.Markdown(fileBytes)),
	}, nil
}
