package sync

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/radovskyb/watcher"
	"github.com/umutsevdi/site/config"
)

type FileCache struct {
	Path string
	Data []byte
}

type Cache struct {
	Static    map[string]FileCache
	Page      map[string]FileCache
	Component map[string]FileCache
}

// Singleton object that stores the caches
var instance Cache

// Whether access to the cache is locked or not
var lock bool

// Whether the cache is initialiced
var ready bool

// An enumerator to access the related cache
const (
	STATIC = iota
	PAGE
	COMPONENT
)

// Returns the cache if exists.
// Waits for the lock when cache is locked
func Get(cache int, key string) *FileCache {
	for lock {

	}
	switch cache {
	case STATIC:
		if data, ok := instance.Static[key]; ok {
			return &data
		}
	case PAGE:
		if data, ok := instance.Page[key]; ok {
			return &data
		}
	case COMPONENT:
		if data, ok := instance.Component[key]; ok {
			return &data
		}
	}
	return nil
}

// Iterates over each item in the related cache
// Waits for the lock when cache is locked
func Each(cache int, function func(k string, data FileCache)) {
	for lock {

	}
	switch cache {
	case STATIC:
		for k, v := range instance.Static {
			function(k, v)
		}
	case PAGE:
		for k, v := range instance.Page {
			function(k, v)
		}
	case COMPONENT:
		for k, v := range instance.Component {
			function(k, v)
		}
	}
}

func init() {
	ready = false
	log.Println("initializing file watcher")
	validateDirectory()
	go startWatcher()
}

// Initializes a watcher that tracks all files under the config#ContentDirectory
func startWatcher() *watcher.Watcher {
	w := watcher.New()
	w.FilterOps(watcher.Rename, watcher.Move, watcher.Write, watcher.Create, watcher.Remove)
	if err := w.AddRecursive(config.ContentDirectory()); err != nil {
		log.Fatalln(err)
	}
	clearMap()
	instance.Static = make(map[string]FileCache, len(w.WatchedFiles()))
	instance.Component = make(map[string]FileCache, len(w.WatchedFiles()))
	instance.Page = make(map[string]FileCache, len(w.WatchedFiles()))

	for k := range w.WatchedFiles() {
		cacheTransactional(strings.Split(k, config.ContentDirectory())[1])
	}
	go onWatchEvent(w)
	log.Println("Watcher started successfully")
	ready = true
	if err := w.Start(time.Second * (time.Duration(config.WatchSecond()))); err != nil {
		log.Fatalln(err)
	}
	return w
}

// A channel handler that updates the cache anytime the file or the directory is updated.
// Restarts the watcher when a file is removed, or added
func onWatchEvent(w *watcher.Watcher) {
	for {
		select {
		case event := <-w.Event:
			if event.IsDir() {
				/* directory events are treated as Write on Linux for some reason */
				log.Println("A directory is updated, restarting the watcher")
				w.Close()
				startWatcher()
			} else {
				/* Update contents of an existing file */
				relPath := strings.Split(event.Path, config.ContentDirectory())[1]
				log.Println(relPath, "is updated")
				cacheTransactional(relPath)
			}
		case err := <-w.Error:
			log.Fatalln(err)
		case <-w.Closed:
			return
		}
	}
}

// Clears the cache
func clearMap() {
	for k := range instance.Static {
		delete(instance.Static, k)
	}
	for k := range instance.Component {
		delete(instance.Component, k)
	}
	for k := range instance.Page {
		delete(instance.Page, k)
	}
}

// Locks the cache, reloads all items from the file system.
func cacheTransactional(path string) {
	lock = true
	if strings.Index(path, "static") == 1 {
		cacheFile(path, config.ContentDirectory()+path, instance.Static)
	} else if strings.Index(path, "page") == 1 {
		ext := filepath.Ext(path)
		if ext == "" {
			return
		}
		key := strings.Split(strings.Split(path, ext)[0], "/page")[1]
		if key == "/index" {
			key = "/"
		}
		cacheFile(key, config.ContentDirectory()+path, instance.Page)
	} else if strings.Index(path, "comp") == 1 {
		key := strings.Split(path, filepath.Ext(path))[0]
		cacheFile(key, config.ContentDirectory()+path, instance.Component)
		wrapComponent(key)
	}
	lock = false
}

// Caches a single file to the desired map
func cacheFile(key string, path string, mapPtr map[string]FileCache) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	mapPtr[key] = FileCache{Data: data, Path: path}
}

// Is map initialized
func IsReady() bool { return ready }
