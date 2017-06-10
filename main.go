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

	"github.com/rjeczalik/notify"
	"github.com/labstack/echo"
	"github.com/olahol/melody"
	"github.com/shurcooL/github_flavored_markdown"
	"time"
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

	watcher := make(chan notify.EventInfo, 1)

	currentFile = getLastUpdated(dirPath)

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

	go func() {
		for {
			select {
			case event := <-watcher:
				log.Println("event:", event)
				if ei, ee := os.Stat(event.Path()); ee != nil || (ei != nil && ei.IsDir()) {
					continue
				}
				if !filepathIsMarkdown(event.Path()) {
					log.Println("not markdown file, continuing...")
					continue
				}
				f := getFilePathFromParam(event.Path())
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
			}
		}
	}()

	if err := notify.Watch(filepath.Join(dirPath, "..."), watcher, notify.All); err != nil {
		log.Fatal(err)
	}
	defer notify.Stop(watcher)

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
		p := c.Param("filename")
		if filepathIsResource(p) {
			return c.File(filepath.Join(dirPath, p))
		}
		filename := getFilePathFromParam(p)
		setCurrentFile(filename)
		// It is important with all this redirecting to NOT allow cacheing.
		c.Response().Header().Set("Cache-Control: no-cache", "true")
		return c.Redirect(http.StatusMovedPermanently, "/")
	})

	log.Println("Listening...", port)
	r.Logger.Fatal(r.Start(":" + strconv.Itoa(port)))
}

func filepathIsMarkdown(path string) bool {
	ff := filepath.Ext(path)
	return !(ff != "" && ff != ".md" && ff != ".markdown" && ff != ".mdown" && ff != ".adoc" && ff != ".txt")
}

func getLastUpdated(path string) (filename string) {
	fs, fe := ioutil.ReadDir(path)
	if fe != nil {
		log.Println(fe)
		return ""
	}
	var latestMod time.Time
	var latestModFile string
	var found bool
	for _, ff := range fs {
		if ff.IsDir() {
			continue
		}
		if !filepathIsMarkdown(ff.Name()) {
			continue
		}
		if ff.ModTime().After(latestMod) {
			found = true
			latestMod = ff.ModTime()
			latestModFile = ff.Name()
		}
	}
	if !found {
		latestModFile = fs[0].Name()
	}
	return latestModFile
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

func filepathIsResource(path string) bool {
	ext := filepath.Ext(path)
	return ext == ".png" || ext == ".jpg" || ext == ".jpeg" || ext == ".svg" || ext == ".tiff" || ext == ".gif"
}

func getFilePathFromParam(param string) string {
	filename := param
	if filename == "" {
		return ""
	}
	if !filepath.IsAbs(filename) {
		filename = filepath.Join(dirPath, filename)
	}
	if fi, e := os.Stat(filename); e == nil && !fi.IsDir() {
		return filename
	}
	if ext := filepath.Ext(filename); ext != "" {
		return filename
	}

	for _, ext := range []string{".md", ".markdown", ".mdown", ".adoc", ".txt"} {
		fname := filename + ext
		if i, e := os.Stat(fname); e == nil && !i.IsDir() {
			return fname
		}
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
