package main

import (
	"errors"
	"mime"
	"net/http"
	"server/syslog"
	"server/util"
	"strconv"
	"strings"
	"time"

	"server/client"
)

func Serve(w http.ResponseWriter, r *http.Request) {
	ip := strings.Split(r.Header.Get("X-Forwarded-For"), ",")[0]
	if ip == "" {
		ip = r.RemoteAddr
	}
	var data []byte
	var err error
	var status int = 200

	switch r.URL.Path {
	case "/":
		data, err = index()
	case "/resume":
		data, err = simplePage("/resume.html")
	case "/robots.txt":
		data, err = util.ReadByteFrom("/robots.txt")
	case "/favicon.ico":
		data, err = util.ReadByteFrom("/favicon.ico")
	case "/sitemap.xml":
		data = []byte(sitemap())
	default:
		err = errors.New("Not found")
	}
	if err != nil {
		data, err = simplePage("/not-found.html")
		status = 404
	}
	data = []byte(setYear(string(data)))
	// Get mime type value and insert it into content header if it exists
	contentHeader := mime.TypeByExtension(util.Ext(r.URL.Path))
	if contentHeader != "" {
		w.Header().Add("Content-Type", contentHeader)
	}
	syslog.Info("main    ", ip, "GET        ", status, r.URL.Path, " ")
	w.WriteHeader(status)
	w.Write(data)
}

func ServeStatic(w http.ResponseWriter, r *http.Request) {
	ip := strings.Split(r.Header.Get("X-Forwarded-For"), ",")[0]
	if ip == "" {
		ip = r.RemoteAddr
	}
	var status int = 200
	data, err := util.ReadByteFrom(r.URL.Path)
	if err != nil {
		status = 404
	}
	// Get mime type value and insert it into content header if it exists
	contentHeader := mime.TypeByExtension(util.Ext(r.URL.Path))
	if contentHeader != "" {
		w.Header().Add("Content-Type", contentHeader)
	}
	syslog.Info("main    ", ip, "GET.static ", status, r.URL.Path[0:], " ")
	w.WriteHeader(status)
	w.Header().Set("Cache-Control", "max-age=3600")
	w.Write(data)
}

func index() ([]byte, error) {
	data, err := util.ReadByteFrom("/index.html")
	var s strings.Builder = strings.Builder{}
	repos, err := client.GetPinnedRepositories(util.Config.User)
	if err != nil {
		syslog.Warn("main     .ReplaceSnippet:", err.Error())
		return data, err
	}
	for _, i := range *repos {
		s.WriteString("<li><div class=\"project-card\"><a target=\"_blank\" href=\"" +
			i.Url + "\"><h2>" + i.Name + "</h2><p class=\"project-icon-block\">")
		for _, j := range i.Language {
			s.WriteString("<span class=\"mono-icon\" style=\"color: " + j.Hex + ";\"> " + client.UnicodeForLanguage(j.Name) + " </span>")
		}
		s.WriteString("<div class=\"project-p-split\"><p>" + i.Description + "</p>")
		s.WriteString("</div></a></div></li>")
	}
	sdata := setYear(string(data))
	return []byte(strings.Replace(sdata, "$PROJECTS", s.String(), 1)), nil
}

func simplePage(url string) ([]byte, error) {
	data, err := util.ReadByteFrom(url)
	if err != nil {
		return nil, err
	}
	return []byte(setYear(string(data))), nil
}

func sitemap() string {
	var s strings.Builder = strings.Builder{}
	s.WriteString(
		"<?xml version=\"1.0\" encoding=\"UTF-8\"?> <urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\">")
	s.WriteString("<url><loc>" + util.Config.Site + "/</loc><priority>1.0</priority></url>")
	routes := []string{"/resume"}
	for _, v := range routes {
		if !strings.Contains(v, "static") && len(v) > 0 && v != "/not-found.html" {
			s.WriteString("<url><loc>" + util.Config.Site + util.File(v) + "</loc> <priority>0.8</priority> </url>")
		}
	}
	s.WriteString("</urlset>")
	return s.String()
}

func setYear(input string) string {
	return strings.Replace(input, "$YEAR", strconv.Itoa(time.Now().Year()), 1)
}
