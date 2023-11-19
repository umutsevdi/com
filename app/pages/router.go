package pages

import (
	"log"
	"mime"
	"net/http"
	"strings"
	"umutsevdi/com/config"

	"github.com/labstack/echo/v4"
)

// Contains ready to render page templates, after being merged with components
var pageCaches map[string][]byte

func Serve(e *echo.Echo) {
	config.StartIndexing()

	config.MapEach("pages", func(k string, v config.FData) {
		e.GET(k, ServePage)
	})
	config.MapEach("static", func(k string, v config.FData) {
		e.GET(k, ServeStatic)
	})
	setGeneratedPages(e)

}

func ServePage(c echo.Context) error {
	ip := c.RealIP()
	path := c.Request().URL.Path
	var file []byte
	var status int = 200

	if fData := config.MapGet("pages", path); fData != nil {
		file = fData.Content
	}
	if file == nil {
		return c.NoContent(http.StatusNotFound)
	}

	contentHeader := mime.TypeByExtension(config.Ext(path))
	if contentHeader == "" {
		contentHeader = "text/html"
	}
	log.Println("GET:", ip, path, status, contentHeader)
	return c.Blob(http.StatusOK, contentHeader, file)
}

func ServeStatic(c echo.Context) error {
	ip := c.RealIP()
	path := c.Request().URL.Path
	var file []byte
	if fData := config.MapGet("static", path); fData != nil {
		file = fData.Content
	} else {
		return c.NoContent(http.StatusNotFound)
	}
	c.Response().Header().Add("Cache-Control", "max-age=3600")

	// Get mime type value and insert it into content header if it exists
	mime := mime.TypeByExtension(config.Ext(path))
	log.Println("GET:", ip, path, mime)
	if mime == "" {
		return c.Blob(http.StatusOK, "text/html", file)
	}
	return c.Blob(http.StatusOK, mime, file)

}

func setGeneratedPages(e *echo.Echo) {
	e.GET("/robots.txt", func(c echo.Context) error {
		path := c.Request().URL.Path
		var file []byte
		if fData := config.MapGet("static", "/static/robots.txt"); fData != nil {
			file = fData.Content
		}
		log.Println("GET:", c.RealIP(), path, http.StatusOK)
		return c.Blob(http.StatusOK, mime.TypeByExtension(config.Ext(path)), file)
	})
	e.GET("/favicon.ico", func(c echo.Context) error {
		path := c.Request().URL.Path
		var file []byte
		if fData := config.MapGet("static", "/static/favicon.ico"); fData != nil {
			file = fData.Content
		}
		log.Println("GET:", c.RealIP(), path, http.StatusOK)
		return c.Blob(http.StatusOK, mime.TypeByExtension(config.Ext(path)), file)
	})
	e.GET("/sitemap.xml", func(c echo.Context) error {
		path := c.Request().URL.Path
		file := []byte(sitemap())
		log.Println("GET:", c.RealIP(), path, http.StatusOK)
		return c.Blob(http.StatusOK, mime.TypeByExtension(config.Ext(path)), file)
	})
	e.RouteNotFound("/*", func(c echo.Context) error {
		return c.HTML(http.StatusNotFound, string(config.MapGet("pages", "/not-found").Content))
	})

	e.RouteNotFound("/static/*", func(c echo.Context) error {
		return c.NoContent(http.StatusNotFound)
	})
}

func sitemap() string {
	var s strings.Builder = strings.Builder{}
	s.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?> <urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\"><url><loc>")
	s.WriteString(*config.C.URI + "/</loc><priority>1.0</priority></url>")

	config.MapEach(
		"pages",
		func(p string, v config.FData) {
			if len(p) > 0 && p != "/not-found" && p != "/" {
				s.WriteString("<url><loc>" + *config.C.URI + p + "</loc> <priority>0.8</priority> </url>")
			}
		})

	s.WriteString("</urlset>")
	return s.String()
}
