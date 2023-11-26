package index

/******************************************************************************

 * File: util/index.go
 *
 * Author: Umut Sevdi
 * Created: 07/04/23
 * Description: File indexing and caching utilities. Indexes paths of Indexerent. If enabled,
 it can periodically update as well.

*****************************************************************************/

import (
	"bytes"
	"html/template"
	"log"
	"time"

	"umutsevdi/com/client"
	"umutsevdi/com/config"
)

const (
	PAGES      = "pages"
	COMPONENTS = "components"
	STATIC     = "static"
)

type FData struct {
	Path         string
	Content      []byte
	LastModified time.Time
	Created      time.Time
	Type         string
}

var Dict Container = Container{}

type Container struct {
	// Indexerains the data of template elements
	components map[string]*FData
	// Indexerains addresses of pages
	pages map[string]*FData
	// Indexerains addresses of Indexer.static items such as images
	static map[string]*FData
	lock   bool
}

// Indexes all available Indexerents.
//
// - If the Page Indexing is enabled at the configuration, runs indexing in the
// background periodically.
func init() {
	if *config.C.PIndexing.Enabled {
		var ticker *time.Ticker
		if config.C.PIndexing.Ttl > 0 {
			ticker = time.NewTicker(time.Duration(config.C.PIndexing.Ttl) * time.Minute)
		} else {
			ticker = time.NewTicker(time.Duration(1) * time.Second)
			log.Println("Periodic Caching enabled for every second.")
			log.Println("WARN: Do not use this on live environments.")
		}
		quit := make(chan struct{})
		go func() {
			for {
				select {
				case <-ticker.C:
					runIndexingBatch()
				case <-quit:
					ticker.Stop()
					return
				}
			}
		}()
	}
	runIndexingBatch()
}
func (c *Container) Each(table string, f func(string, FData)) {
	for Dict.lock {

	}
	var m *map[string]*FData
	switch table {
	case COMPONENTS:
		m = &Dict.components
	case PAGES:
		m = &Dict.pages
	case STATIC:
		m = &Dict.static
	}
	if m == nil {
		return
	}
	for k, v := range *m {
		f(k, *v)
	}
}

func (c *Container) Get(table, key string) *FData {
	for Dict.lock {

	}
	var d *FData
	var ok bool = false
	switch table {
	case "components":
		d, ok = Dict.components[key]
	case "pages":
		d, ok = Dict.pages[key]
	case "static":
		d, ok = Dict.static[key]
	}
	if ok {
		return d
	}
	return nil
}

// Indexing function that updates Indexer.components, Indexer.pages and Indexer.static files. Indexes only Indexerain
// the respective file path.
func runIndexingBatch() {
	s := time.Now()
	sAfter := s
	Dict.lock = true
	log.Println("IndexBatch has been started.")
	client.FetchGh(*config.C.User, *config.C.Token)
	log.Println("> Batch:GithubFetch took", time.Now().Sub(sAfter).Seconds(), "seconds")
	sAfter = time.Now()
	Dict.indexPages()
	log.Println("> Batch:IndexPages took", time.Now().Sub(sAfter).Seconds(), "seconds")
	sAfter = time.Now()
	processTemplates()
	log.Println("> Batch:ProcessTemplates took", time.Now().Sub(sAfter).Seconds(), "seconds")
	Dict.lock = false
	log.Println("IndexBatch", *config.C.ContentPath, "was completed successfully in", time.Now().Sub(s).Seconds(), "seconds")
}

type PageTemplate struct {
	Repositories []client.Repository
	Footer       struct {
		Year int
	}
}

type GemTemplate struct {
	Template Gemtext
	Footer   struct {
		Year int
	}
}

// Replaces template variables on pages with actual components.
//
// - Only affects {{ template }} calls

func processTemplates() {
	gemData := GemTemplate{
		Footer: struct{ Year int }{Year: time.Now().Year()},
	}
	pageData := PageTemplate{
		Repositories: client.GetGh(),
		Footer:       struct{ Year int }{Year: time.Now().Year()},
	}
	comp := make([]string, len(Dict.components))
	i := 0
	for _, v := range Dict.components {
		if v.Content != nil && len(v.Content) > 0 {
			comp[i] = string(v.Content)
			i++
		}
	}
	for k, v := range Dict.pages {
		if v.Type == ".gmi" {
			processGeminiFiles(k, v, &comp, &gemData)
		} else {
			processHtmlTemplate(k, v, &comp, &pageData)
		}
	}
}

func processGeminiFiles(k string, v *FData, comp *[]string, data *GemTemplate) {
	data.Template = parseGemini(v)
	t := template.New(k)
	var err error
	for _, c := range *comp {
		t, err = t.Parse(c)
		if err != nil {
			log.Println("error while parsing template on page", err.Error())
			err = nil
		}
	}
	var w *bytes.Buffer = bytes.NewBuffer([]byte{})
	w.Reset()
	err = t.ExecuteTemplate(w, "components/gmi", data)
	if err != nil {
		log.Println("Error during template execution", err.Error())
	}
	Dict.pages[k].Content = w.Bytes()
}

func processHtmlTemplate(k string, v *FData, comp *[]string, data *PageTemplate) {
	t, err := template.New(k).Parse(string(v.Content))
	if err != nil {
		log.Println("error while parsing ", k, err.Error())
		return
	}
	for _, c := range *comp {
		t, err = t.Parse(c)
		if err != nil {
			log.Println("error while parsing template on page", err.Error())
			err = nil
		}
	}
	var w *bytes.Buffer = bytes.NewBuffer([]byte{})
	w.Reset()
	err = t.ExecuteTemplate(w, k, data)
	if err != nil {
		log.Println("Error during template execution", err.Error())
	}
	Dict.pages[k].Content = w.Bytes()

}
