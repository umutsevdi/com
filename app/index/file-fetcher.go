package index

import (
	"fmt"
	"log"
	"os"
	"strings"
	"umutsevdi/com/config"
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
	for _, v := range data {
		if v.Type().IsDir() {
			indexHtml(basePath, path+v.Name()+"/", uris)
		} else {
			var key string = path + v.Name()
			var path string = *config.C.ContentPath + basePath + path + v.Name()
			key = strings.Split(key, ".")[0]
			if key == "/index" {
				key = "/"
			}
			mapToCache(key, path, uris)
		}
	}
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
		metadata, err := os.Stat(fs.Path)
		if err != nil {
			log.Println("WARN: File", path, "does not exist")
			return
		}
		fs.LastModified = metadata.ModTime()
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
		metadata, err := os.Stat(path)
		if err != nil {
			log.Println("WARN: File", path, "does not exist for metadata")
			return
		}
		// Update only if the file is changed
		(*files)[key] = &FData{
			Path:         path,
			Content:      d,
			LastModified: metadata.ModTime(),
		}
	}
}
