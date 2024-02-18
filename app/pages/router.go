package pages

import (
	"log"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/umutsevdi/site/config"
	"github.com/umutsevdi/site/sync"
)

func Start() {
	e := echo.New()
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{Level: 5}))
	Serve(e)
}

func Serve(e *echo.Echo) {
	for !sync.IsReady() {
	}
	sync.Each(sync.STATIC, func(k string, v sync.FileCache) {
		e.GET(k, ServeStatic)
	})
	sync.Each(sync.PAGE, func(k string, v sync.FileCache) {
		log.Println(k)
		e.GET(k, ServePage)
	})
	MapGeneratedPages(e)
	str := ":" + strconv.Itoa(int(config.Port()))
	log.Fatal(e.Start(str))

}

// Serve a static file under the /static path
func ServeStatic(c echo.Context) error {
	ip := c.RealIP()
	path := c.Request().URL.Path
	content := sync.Get(sync.STATIC, path)
	if content == nil {
		log.Println("GET:", ip, path, http.StatusNotFound)
		return c.NoContent(http.StatusNotFound)
	}
	c.Response().Header().Add("Cache-Control", "max-age=3600")

	// Get mime type value and insert it into content header if it exists
	mime := mime.TypeByExtension(filepath.Ext(path))
	log.Println("GET:", ip, path, mime)
	if mime == "" {
		return c.Blob(http.StatusOK, "text/html", content.Data)
	}
	return c.Blob(http.StatusOK, mime, content.Data)
}

// Serves a regular page
// TODO this function will be replaced with individual routing functions
func ServePage(c echo.Context) error {
	ip := c.RealIP()
	path := c.Request().URL.Path
	if path == "" {
		path = "/"
	}
	content := sync.Get(sync.PAGE, path)
	if content == nil {
		return c.NoContent(http.StatusNotFound)
	}
	log.Println("GET:", ip, path, 200)
	return c.HTMLBlob(http.StatusOK, sync.ProcessTemplates(path, content, sync.Data()))
}
