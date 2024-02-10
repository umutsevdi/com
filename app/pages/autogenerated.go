package pages

import (
	"log"
	"mime"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/umutsevdi/site/config"
	"github.com/umutsevdi/site/sync"
)

// Maps all runtime generated pages that does not correspond to a file
func MapGeneratedPages(e *echo.Echo) {
	e.GET("/robots.txt", func(c echo.Context) error { return routeStaticContent(c, "/static/robots.txt") })
	e.GET("/favicon.ico", func(c echo.Context) error { return routeStaticContent(c, "/static/favicon.ico") })
	e.GET("/sitemap.xml", sitemap)
	e.GET("/sitemaps.xml", sitemap)
	e.RouteNotFound("/*", func(c echo.Context) error {
		key := "/not-found"
		content := sync.Get(sync.PAGE, key)
		return c.HTMLBlob(http.StatusNotFound,
			sync.ProcessTemplates(key, content, sync.Data()))
	})
	e.RouteNotFound("/static/*", func(c echo.Context) error {
		return c.NoContent(http.StatusNotFound)
	})
}

// Proxies a file under the /static directory to its proper URI
func routeStaticContent(c echo.Context, path string) error {
	content := sync.Get(sync.STATIC, path)
	if content == nil {
		return c.NoContent(http.StatusNotFound)
	}
	return c.Blob(http.StatusOK, mime.TypeByExtension(content.Path), content.Data)
}

// Builds a sitemap.xml file and serves it
func sitemap(c echo.Context) error {
	log.Println("GET:", c.RealIP(), c.Request().URL.Path, http.StatusOK)
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

	return c.Blob(http.StatusOK, mime.TypeByExtension(".xml"), []byte(s.String()))
}
