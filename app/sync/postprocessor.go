package sync

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"text/template"
	"time"

	"github.com/umutsevdi/site/client"
	"github.com/umutsevdi/site/config"
)

// Wraps a component with is name to define the template
func wrapComponent(key string) {
	comp := instance.Component[key]
	comp.Data = []byte(fmt.Sprintf("{{define \"%s\"}}\n%s{{end}}",
		key, comp.Data))
	instance.Component[key] = comp
}

// Checks whether the folder structure fulfills  the requirements. Exits on
// error
func validateDirectory() {
	file, err := os.Open(config.ContentDirectory())
	if err != nil {
		log.Fatal("ERROR: File not found", config.ContentDirectory())
	}
	defer file.Close()
	files, err := file.Readdirnames(-1)
	if err != nil {
		log.Fatal("ERROR: File not found", config.ContentDirectory())
	}
	directories := 0
	for i := range files {
		if files[i] == "static" || files[i] == "comp" || files[i] == "page" {
			directories++
		}
	}
	if directories != 3 {
		log.Fatal("ERROR: Folder structure is invalid. Expected directories: page, static, comp")
	}
}

// Processes the template on the given path, inserts the predefined templates and
// fills the data
//
// returns final page as byte array
func ProcessTemplates(path string, content *FileCache, data interface{}) []byte {
	t, err := template.New(path).Parse(string(content.Data))
	if err != nil {
		log.Println("error while parsing template on page", err.Error())
		err = nil
	}
	for _, comp := range instance.Component {
		t, err = t.Parse(string(comp.Data))
		if err != nil {
			log.Println("error while parsing template on page", err.Error())
			err = nil
		}
	}
	var w *bytes.Buffer = bytes.NewBuffer([]byte{})
	w.Reset()
	err = t.ExecuteTemplate(w, path, data)
	return w.Bytes()
}

type PageTemplate struct {
	Repositories []client.Repository
	Footer       struct {
		Year int
	}
}

// Temporary structure and its getter function that stores the data to generate index.html
func Data() PageTemplate {
	return PageTemplate{
		Repositories: client.GetGh(),
		Footer:       struct{ Year int }{Year: time.Now().Year()},
	}
}
