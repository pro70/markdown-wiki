package pages

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"time"

	"github.com/julienschmidt/httprouter"

	"github.com/irgangla/markdown-wiki/log"
	"github.com/irgangla/markdown-wiki/sdk"
	wikiTemplate "github.com/irgangla/markdown-wiki/template"
)

// EditContent contains the edit data
type EditContent struct {
	sdk.MetaData
	Name    string
	Content string
}

// SaveInput for edit json request data
type SaveInput struct {
	Name    string
	Content string
}

// SaveOutput for edit json response data
type SaveOutput struct {
	Result   string
	Text     string
	Time     string
	Duration string
}

// EditEndpoint for markdown pages
func EditEndpoint(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	log.Info("EDIT", "URI", request.RequestURI)

	writer.Header().Set("Content-Type", "text/html")

	name := params.ByName("name")
	path := filepath.Join(".", "data", "md", name+".md")
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Error("EDIT", "Read file", err)
		writer.WriteHeader(500)
		writer.Write([]byte(err.Error()))
		return
	}

	var content EditContent
	content.Title = "Edit " + name
	content.Description = "Edit page " + name
	content.Content = string(data)
	content.Name = name

	var layouts = make([]string, 0)
	layouts = append(layouts, "/css/edit.css")
	content.Layouts = layouts

	var scripts = make([]string, 0)
	scripts = append(scripts, "/js/edit.js")
	content.Scripts = scripts

	t, err := wikiTemplate.Get("edit.html")
	if err != nil || t == nil {
		log.Error("EDIT", "Template error")
		writer.WriteHeader(500)
		if err != nil {
			log.Error("EDIT", "Template error", err)
			writer.Write([]byte(err.Error()))
		}
		return
	}

	writer.WriteHeader(200)

	err = t.Execute(writer, content)
	if err != nil {
		log.Error("EDIT", "Template render error", err)
		return
	}
}

// SaveEndpoint for update markdown pages
func SaveEndpoint(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	timer := time.Now()

	log.Info("SAVE", "Save called from", request.RemoteAddr)

	var data SaveInput
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		log.Error("SAVE", "request data error", err)
		var responseData SaveOutput
		responseData.Result = "nok"
		responseData.Text = "problem with user json data"
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(400)
		err = json.NewEncoder(writer).Encode(responseData)
		if err != nil {
			log.Error("SAVE", err)
		}
		return
	}

	log.Info("SAVE", data.Name)

	err = updateFile(data.Name, data.Content)
	if err != nil {
		log.Error("SAVE", "update file error", err)
		var responseData SaveOutput
		responseData.Result = "nok"
		responseData.Text = "update file failed"
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(500)
		err = json.NewEncoder(writer).Encode(responseData)
		if err != nil {
			log.Error("SAVE", err)
		}
		return
	}

	end := time.Since(timer)
	log.Info("SAVE", "Processing takes:", end.String())

	pageUpdateEvent(data.Name)

	var responseData SaveOutput
	responseData.Result = "ok"
	responseData.Text = data.Name + " updated"
	responseData.Time = time.Now().Format("02/01/2006, 15:04:05")
	responseData.Duration = end.String()
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(200)
	err = json.NewEncoder(writer).Encode(responseData)
	if err != nil {
		log.Error("SAVE", "response encode error", err)
	}
	return
}

func pageUpdateEvent(name string) {
	var event sdk.Event
	event.Event = "MD_UPDATE"
	event.Data = name
	sdk.ClientEvents <- event
}

func updateFile(name string, content string) error {
	path := filepath.Join("data", "md", name+".md")
	log.Info("SAVE", "Update file", path)
	return ioutil.WriteFile(path, []byte(content), 0775)
}
