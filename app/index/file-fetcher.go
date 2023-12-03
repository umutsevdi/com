package index

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
	"umutsevdi/com/config"

	"github.com/djherbis/times"
)

func (c *Container) indexPages() {
	if len(c.components) == 0 {
		c.components = make(map[string]*FData, 100)
	}
	if len(c.pages) == 0 {
		c.pages = make(map[string]*FData, 100)
	}
	if len(c.static) == 0 {
		c.static = make(map[string]*FData, 100)
	}
	indexHtml("/components", "/", &c.components)
	indexHtml("/pages", "/", &c.pages)
	c.indexStatic("/")
	c.defineComponents()
}

// Recursively traverses through the content/pages directory
// and indexes pages.
//
// - File extensions are removed while indexing
//
//	@param path - path to traverse
//	@param uris - map to insert items
//
// "content/pages"
// "content/components"
func indexHtml(basePath string, path string, uris *map[string]*FData) {
	data, err := os.ReadDir(*config.C.ContentPath + basePath + path)
	if err != nil {
		return
	}
	var dKeys [][2]string
	for _, v := range data {
		if v.Type().IsDir() {
			dKeys = append(dKeys, [2]string{path + v.Name() + "/", v.Name() + "/"})
			indexHtml(basePath, path+v.Name()+"/", uris)
		} else {
			var key string = path + v.Name()
			var path string = *config.C.ContentPath + basePath + path + v.Name()
			key = strings.Split(key, ".")[0]
			if key == "/index" {
				key = "/"
			}
			mapToCache(key, path, uris)
			title := strings.Split(key, "/")
			dKeys = append(dKeys, [2]string{key, title[len(title)-1]})
		}
	}
	mapDirToCache(path, dKeys)
}

// Recursively traverses through the content/static directory
// and indexes static files at /static URL path.
//
//	@param path - path to traverse
//	@param uris - map to insert items
func (c *Container) indexStatic(path string) {
	data, err := os.ReadDir(*config.C.ContentPath + "/static" + path)
	if err != nil {
		return
	}
	for _, v := range data {
		if v.Type().IsDir() {
			c.indexStatic(path + v.Name() + "/")
		} else {
			var key string = "/static" + path + v.Name()
			var path string = *config.C.ContentPath + key
			mapToCache(key, path, &c.static)
		}
	}
}

// Inserts the name of the template to be resolved by the engine
func (c *Container) defineComponents() {
	for k, v := range c.components {
		v.Content = []byte(fmt.Sprintf("{{define \"components%s\"}}\n%s{{end}}",
			k, c.components[k].Content))
	}
}

// Caches given path for the key
// - If the file is not registered, inserts it with last modified date
// - If the file is already cached, checks whether it has been updated or not
// updates only it's changed
func mapToCache(key, path string, files *map[string]*FData) {
	if fs, ok := (*files)[key]; ok {
		t, err := times.Stat(fs.Path)
		if err != nil {
			log.Println("WARN: File", path, "does not exist")
			return
		}
		if t.HasChangeTime() {
			fs.LastModified = t.ChangeTime()
		}
		if t.HasBirthTime() {
			fs.Created = t.BirthTime()
		}
		fs.Content, err = os.ReadFile(fs.Path)
		if err != nil {
			log.Println("WARN: File", path, "does not exist for metadata")
			return
		}
	} else {
		// if doesn't exist create
		d, err := os.ReadFile(path)
		if err != nil {
			log.Println("WARN: File", path, "does not exist")
			return
		}
		t, err := times.Stat(path)
		if err != nil {
			log.Println("WARN: File", path, "does not exist for metadata")
			return
		}
		fs := FData{
			Path:    path,
			Content: d,
			Type:    ext(path),
		}
		if t.HasChangeTime() {
			fs.LastModified = t.ChangeTime()
		}
		if t.HasBirthTime() {
			fs.Created = t.BirthTime()
		}
		// Update only if the file is changed
		(*files)[key] = &fs
	}
}

type DirTemplate struct {
	Key  string
	Uris [][2]string
}

func mapDirToCache(key string, uris [][2]string) {
	if key == "/" {
		return
	}
	var templateString = `// Title: {{.Key}}
// Author: Umut Sevdi
{{range .Uris}}
=>{{index . 0}} {{index . 1}}
{{end}}
`
	var w *bytes.Buffer = bytes.NewBuffer([]byte{})
	t, err := template.New("dir").Parse(templateString)

	if err != nil {
		log.Println("WARN: Dir index error", err)
	}
	t.Execute(w, DirTemplate{Key: key, Uris: uris})
	data := FData{
		Content: w.Bytes(),
		Path:    key,
		Type:    "/.gmi",
	}
	fmt.Println(key, string(data.Content))

	Dict.pages[key] = &data
}

// Parses received URL and extracts it's extension
//
//	@param URL path to file
//	@return string corresponding extension type in the format of .type
func ext(url string) string {
	p := strings.Split(url, ".")
	return "." + p[len(p)-1]
}
