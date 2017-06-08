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
	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
	"github.com/olahol/melody"
	"github.com/rjeczalik/notify"
	"github.com/shurcooL/github_flavored_markdown"
)

var port int
var dirPath string
var ei chan notify.EventInfo

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
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)

					f := getFilePathFromParam(event.Name)
					m, e := getReadFile(f)
					if e != nil {
						log.Println(e)
						return
					}
					b, e := json.Marshal(m)
					if e != nil {
						log.Println(e)
						return
					}
					mm.Broadcast(b)
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(dirPath)
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "index.html")
	})
	r.GET("/0", func(c *gin.Context) {
		mm.HandleRequest(c.Writer, c.Request)
	})

	log.Println("Listening...", port)
	r.Run(":" + strconv.Itoa(port))

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
	return filename
}

func deliverMarkdownedFileFromParams(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	f := getFilePathFromParam(params["filename"])
	m, e := getReadFile(f)
	if e != nil {
		log.Println(e)
		return
	}
	b, e := json.Marshal(m)
	if e != nil {
		log.Println(e)
		return
	}
	w.Write(b)
}

func getReadFile(path string) (FileContent, error) {
	fileBytes, e := ioutil.ReadFile(path)
	if e != nil {
		log.Println(e)
		return FileContent{}, e
	}
	return FileContent{
		Title: filepath.Base(path), // TODO parse File-Name.md syntax => File Name
		Body:  string(github_flavored_markdown.Markdown(fileBytes)),
	}, nil
}
