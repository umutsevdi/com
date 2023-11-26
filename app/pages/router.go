package pages

import (
	"log"
	"mime"
	"net/http"
	"strings"
	"umutsevdi/com/config"
	"umutsevdi/com/index"

	"github.com/labstack/echo/v4"
)

func Serve(e *echo.Echo) {
	index.Dict.Each(index.PAGES, func(k string, v index.FData) {
		e.GET(k, ServePage)
	})
	index.Dict.Each(index.STATIC, func(k string, v index.FData) {
		e.GET(k, ServeStatic)
	})
	setGeneratedPages(e)

}

func ServePage(c echo.Context) error {
	ip := c.RealIP()
	path := c.Request().URL.Path
	if path == "" {
		path = "/"
	}
	fData := index.Dict.Get(index.PAGES, path)
	if fData == nil {
		return c.NoContent(http.StatusNotFound)
	}
	fType := mime.TypeByExtension(fData.Type)
	if fType == "" {
		fType = "text/html"
	}

	log.Println("GET:", ip, path, 200, fType)
	return c.Blob(http.StatusOK, fType, fData.Content)
}

func ServeStatic(c echo.Context) error {
	ip := c.RealIP()
	path := c.Request().URL.Path
	fData := index.Dict.Get("static", path)
	if fData == nil {
		return c.NoContent(http.StatusNotFound)
	}
	c.Response().Header().Add("Cache-Control", "max-age=3600")

	// Get mime type value and insert it into content header if it exists
	mime := mime.TypeByExtension(fData.Type)
	log.Println("GET:", ip, path, mime)
	if mime == "" {
		return c.Blob(http.StatusOK, "text/html", fData.Content)
	}
	return c.Blob(http.StatusOK, mime, fData.Content)

}

func setGeneratedPages(e *echo.Echo) {
	e.GET("/robots.txt", func(c echo.Context) error {
		path := c.Request().URL.Path
		fData := index.Dict.Get(index.STATIC, "/static/robots.txt")
		if fData == nil {
			return c.NoContent(http.StatusNotFound)
		}
		log.Println("GET:", c.RealIP(), path, http.StatusOK)
		return c.Blob(http.StatusOK, mime.TypeByExtension(fData.Type), fData.Content)
	})
	e.GET("/favicon.ico", func(c echo.Context) error {
		path := c.Request().URL.Path
		fData := index.Dict.Get(index.STATIC, "/static/favicon.ico")
		if fData == nil {
			return c.NoContent(http.StatusNotFound)
		}
		log.Println("GET:", c.RealIP(), path, http.StatusOK)
		return c.Blob(http.StatusOK, mime.TypeByExtension(fData.Type), fData.Content)
	})
	e.GET("/sitemap.xml", func(c echo.Context) error {
		path := c.Request().URL.Path
		file := []byte(sitemap())
		log.Println("GET:", c.RealIP(), path, http.StatusOK)
		return c.Blob(http.StatusOK, mime.TypeByExtension(".xml"), file)
	})
	e.RouteNotFound("/*", func(c echo.Context) error {
		return c.HTML(http.StatusNotFound, string(index.Dict.Get(index.PAGES, "/not-found").Content))
	})

	e.RouteNotFound("/static/*", func(c echo.Context) error {
		return c.NoContent(http.StatusNotFound)
	})
}

func sitemap() string {
	var s strings.Builder = strings.Builder{}
	s.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?> <urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\"><url><loc>")
	s.WriteString(*config.C.URI + "/</loc><priority>1.0</priority></url>")
	index.Dict.Each(index.PAGES, func(p string, v index.FData) {
		if len(p) > 0 && p != "/not-found" && p != "/" {
			s.WriteString("<url><loc>" + *config.C.URI + p + "</loc> <priority>0.8</priority> </url>")
		}
	})
	s.WriteString("</urlset>")
	return s.String()
}
