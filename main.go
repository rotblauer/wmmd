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
)

var port int
var dirPath string

type FileContent struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

func init() {
	flag.IntVar(&port, "port", 3000, "port to serve on")
}

func main() {
	flag.Parse()
	dirPath = mustMakeDirPath()
	mm := melody.New()

	var currentFile string = filepath.Join(dirPath, "README.md")

	mm.HandleConnect(func(s *melody.Session) {
		log.Println("session connected")
		sidebar, e := getReadFile(filepath.Join(dirPath, "_Sidebar.md"))
		if e != nil {
			log.Println(e)
		}
		footer, e := getReadFile(filepath.Join(dirPath, "_Footer.md"))
		if e != nil {
			log.Println(e)
		}
		curFile, e := getReadFile(currentFile)
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
				f := getFilePathFromParam(event.Name)
				currentFile = f
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
				log.Println("error:", err)
			}
		}
	}()
	// Sometimes on refresh watcher stops silently.
	err = watcher.Add(dirPath)
	if err != nil {
		log.Fatal(err)
	}

	r := echo.New()
	r.File("/", "index.html")
	// Static assets.
	r.Static("/assets", "./assets")
	r.Static("/node_modules/primer-css/build", "./node_modules/primer-css/build")
	// Websocket.
	r.GET("/x/0", func(c echo.Context) error {
		mm.HandleRequest(c.Response(), c.Request())
		return nil
	})
	// Any other filename.
	r.Any("/:filename", func(c echo.Context) error {
		filename := c.Param("filename")
		filename = getFilePathFromParam(filename)
		currentFile = filename
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
