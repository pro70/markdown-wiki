package template

import (
	"html/template"
	"io/ioutil"
	"path/filepath"

	"github.com/irgangla/markdown-wiki/log"
)

var (
	sharedTemplates []string
	templates       map[string]*template.Template
)

// Get template with given name
func Get(name string) (*template.Template, error) {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}
	template, ok := templates[name]
	if !ok {
		files := make([]string, 0)

		file := filepath.Join(".", "data", "template", name)
		files = append(files, file)

		shared := shared()
		files = append(files, shared...)

		var err error
		template, err = template.ParseFiles(files...)
		if err != nil {
			log.Error("TEMPLATE", "invalid template", err)
			return nil, err
		}
		templates[name] = template
	}
	return template, nil
}

func shared() []string {
	if sharedTemplates == nil {
		dir := filepath.Join(".", "data", "template", "shared")
		infos, err := ioutil.ReadDir(dir)
		if err == nil {
			sharedTemplates = make([]string, 0)
			for _, i := range infos {
				file := filepath.Join(dir, i.Name())
				sharedTemplates = append(sharedTemplates, file)
			}
		} else {
			log.Error("TEMPLATE", "find shared templates", err)
		}
	}
	return sharedTemplates
}
