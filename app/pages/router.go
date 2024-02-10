package pages

import (
	"log"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

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
	setGeneratedPages(e)
	str := ":" + strconv.Itoa(int(config.Port()))
	log.Fatal(e.Start(str))

}

func ServeStatic(c echo.Context) error {
	ip := c.RealIP()
	path := c.Request().URL.Path
	log.Println(path)
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
	mime := "text/html"
	log.Println("GET:", ip, path, 200, mime)
	return c.Blob(http.StatusOK, mime, sync.ResolvePage(path, content))
}

func setGeneratedPages(e *echo.Echo) {
	e.GET("/robots.txt", func(c echo.Context) error {
		path := c.Request().URL.Path
		content := sync.Get(sync.STATIC, "/static/robots.txt")
		if content == nil {
			return c.NoContent(http.StatusNotFound)
		}
		log.Println("GET:", c.RealIP(), path, http.StatusOK)
		return c.Blob(http.StatusOK, mime.TypeByExtension(content.Path), content.Data)
	})
	e.GET("/favicon.ico", func(c echo.Context) error {
		path := c.Request().URL.Path
		content := sync.Get(sync.STATIC, "/static/favicon.ico")
		if content == nil {
			return c.NoContent(http.StatusNotFound)
		}
		log.Println("GET:", c.RealIP(), path, http.StatusOK)
		return c.Blob(http.StatusOK, mime.TypeByExtension(content.Path), content.Data)
	})
	e.GET("/sitemap.xml", func(c echo.Context) error {
		path := c.Request().URL.Path
		file := []byte(sitemap())
		log.Println("GET:", c.RealIP(), path, http.StatusOK)
		return c.Blob(http.StatusOK, mime.TypeByExtension(".xml"), file)
	})
	e.RouteNotFound("/*", func(c echo.Context) error {
		return c.HTML(http.StatusNotFound, string(sync.Get(sync.PAGE, "/not-found").Data))
	})
	e.RouteNotFound("/static/*", func(c echo.Context) error {
		return c.NoContent(http.StatusNotFound)
	})
}

func sitemap() string {
	var s strings.Builder = strings.Builder{}
	s.WriteString(`<?xml version="1.0" encoding="UTF-8"?>
    <urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
    <url><loc>`)
	s.WriteString(config.URI() + "/</loc><priority>1.0</priority></url>")
	sync.Each(sync.PAGE, func(p string, v sync.FileCache) {
		if len(p) > 0 && p != "/not-found" && p != "/" {
			s.WriteString("<url><loc>" + config.URI() + p +
				"</loc> <priority>0.8</priority> </url>")
		}
	})
	s.WriteString("</urlset>")
	return s.String()
}
