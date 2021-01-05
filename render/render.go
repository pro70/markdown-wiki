package render

import (
	"bytes"
	"errors"
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"

	chromaHtml "github.com/alecthomas/chroma/formatters/html"
	"github.com/julienschmidt/httprouter"
	mathjax "github.com/litao91/goldmark-mathjax"
	"github.com/yuin/goldmark"
	emoji "github.com/yuin/goldmark-emoji"
	highlighting "github.com/yuin/goldmark-highlighting"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"

	"github.com/irgangla/markdown-wiki/log"
	wikiTemplate "github.com/irgangla/markdown-wiki/template"
)

// MetaDataContent contains the metadata of the rendered file
type MetaDataContent struct {
	Title       string
	Description string
	Author      string
	Layouts     []string
	Scripts     []string
}

// MarkdownContent contains the render result
type MarkdownContent struct {
	MetaDataContent
	Content template.HTML
}

// Endpoint to server markdown pages
func Endpoint(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	log.Info("MARKDOWN", "URI", request.RequestURI)

	writer.Header().Set("Content-Type", "text/html")

	data, err := renderMarkdown(params.ByName("name"))
	if err != nil || data == nil {
		log.Error("MARKDOWN", "Data error")
		writer.WriteHeader(404)
		if err != nil {
			log.Error("MARKDOWN", "Data error", err)
			writer.Write([]byte(err.Error()))
		}
		return
	}

	t, err := wikiTemplate.Get("markdown.html")
	if err != nil || t == nil {
		log.Error("MARKDOWN", "Template error")
		writer.WriteHeader(500)
		if err != nil {
			log.Error("MARKDOWN", "Template error", err)
			writer.Write([]byte(err.Error()))
		}
		return
	}

	writer.WriteHeader(200)

	data.Scripts = append(data.Scripts, "https://polyfill.io/v3/polyfill.min.js?features=es6")
	data.Scripts = append(data.Scripts, "https://cdn.jsdelivr.net/npm/mathjax@3/es5/tex-mml-chtml.js")

	err = t.Execute(writer, *data)
	if err != nil {
		log.Error("MARKDOWN", "Template render error", err)
		return
	}
}

// File renders the given markdown file as HTML string
func File(name string) (*MarkdownContent, error) {
	data, err := renderMarkdown(name)
	if err != nil || data == nil {
		log.Error("MARKDOWN", "File render error")
		if err != nil {
			log.Error("MARKDOWN", "File render error", err)
			return nil, err
		}
		return nil, errors.New("file render error")
	}
	return data, nil
}

func renderMarkdown(name string) (*MarkdownContent, error) {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Footnote,
			extension.DefinitionList,
			meta.Meta,
			highlighting.NewHighlighting(
				highlighting.WithStyle("monokai"),
				highlighting.WithFormatOptions(
					chromaHtml.WithLineNumbers(true),
				),
			),
			emoji.Emoji,
			mathjax.MathJax,
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
		),
	)
	var buf bytes.Buffer
	path := filepath.Join(".", "data", "md", name+".md")
	source, err := ioutil.ReadFile(path)
	if err != nil {
		log.Error("MARKDOWN", "markdown parsing error", err)
		return nil, err
	}
	context := parser.NewContext()
	err = md.Convert(source, &buf, parser.WithContext(context))
	if err != nil {
		log.Error("MARKDOWN", "markdown parsing error", err)
		return nil, err
	}

	metaData := meta.Get(context)
	data := new(MarkdownContent)

	data.Title = getMetaDataString("Title", &metaData)
	data.Description = getMetaDataString("Summary", &metaData)
	data.Author = getMetaDataString("Author", &metaData)
	data.Layouts = getMetaDataList("Layouts", &metaData)
	data.Scripts = getMetaDataList("Scripts", &metaData)

	data.Content = template.HTML(buf.String())

	return data, nil
}

func getMetaDataString(key string, data *map[string]interface{}) string {
	d, ok := (*data)[key]
	if !ok {
		return ""
	}
	v, ok := d.(string)
	if !ok {
		return ""
	}
	return v
}

func getMetaDataList(key string, data *map[string]interface{}) []string {
	d, ok := (*data)[key]
	if !ok {
		return make([]string, 0)
	}
	v, ok := d.([]string)
	if !ok {
		return make([]string, 0)
	}
	return v
}
