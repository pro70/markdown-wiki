package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/irgangla/markdown-wiki/log"
)

func initializeData() {
	makeDirs()

	createIfNotExist(filepath.Join("data", "css", "layout.css"), layoutCSS)
	createIfNotExist(filepath.Join("data", "template", "shared", "layout.html"), layoutHTML)
	createIfNotExist(filepath.Join("data", "template", "markdown.html"), markdownHTML)
	createIfNotExist(filepath.Join("data", "css", "edit.css"), editCSS)
	createIfNotExist(filepath.Join("data", "js", "edit.js"), editJS)
	createIfNotExist(filepath.Join("data", "template", "edit.html"), editHTML)

	createEmojiList()

	copyInitData()
}

func copyInitData() {
	err := filepath.Walk("initData",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Error("INIT", "copy init data", err)
				return err
			}

			if info.IsDir() {
				dir := "data" + path[8:]
				log.Info("INIT", "make dir", dir)
				os.MkdirAll(dir, 0775)
				return nil
			}

			destinationFile := "data" + path[8:]
			return copyFileIfNotExists(path, destinationFile)
		})
	if err != nil {
		log.Error("INIT", "walk init data", err)
	}
}

func copyFileIfNotExists(sourceFile, destinationFile string) error {
	if _, err := os.Stat(destinationFile); os.IsNotExist(err) {
		log.Info("INIT", "copy init file", sourceFile)
		input, err := ioutil.ReadFile(sourceFile)
		if err != nil {
			log.Error("INIT", "read source", sourceFile, err)
			return err
		}
		err = ioutil.WriteFile(destinationFile, input, 0775)
		if err != nil {
			log.Error("INIT", "write file", destinationFile, err)
			return err
		}
	}
	return nil
}

func makeDirs() {
	dirs := [][]string{
		[]string{".", "data", "css"},
		[]string{".", "data", "html"},
		[]string{".", "data", "js"},
		[]string{".", "data", "md"},
		[]string{".", "data", "template"},
		[]string{".", "data", "template", "shared"},
	}

	for _, d := range dirs {
		path := filepath.Join(d...)
		err := os.MkdirAll(path, 0775)
		if err != nil {
			log.Error("INIT", "Make dir failed:", d, err)
		}
	}
}

func createIfNotExist(path string, content string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := ioutil.WriteFile(path, []byte(content), 0775)
		if err != nil {
			log.Error("INIT", "Write file failed:", path, err)
		}
	}
}

func createEmojiList() {
	path := filepath.Join("data", "md", "emoji.md")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		content := emojiMeta + "\n\n# Emoji\n\n"
		data := getEmoji()
		for k, vals := range data {
			content = content + "## " + k + "\n\n"
			for _, v := range vals {
				content = content + "* `" + v + "` " + v + "\n"
			}
		}
		err := ioutil.WriteFile(path, []byte(content), 0775)
		if err != nil {
			log.Error("INIT", "Write file failed:", path, err)
		}
	}
}

func getEmoji() map[string][]string {
	emojis := make(map[string][]string)

	data, err := ioutil.ReadFile(filepath.Join("initData", "emoji"))
	if err != nil {
		log.Error("INIT", "load emoji", err)
		return emojis
	}

	content := string(data)
	lines := strings.Split(content, "\n")
	for _, l := range lines {
		parts := strings.Split(l, " ")
		name := parts[0]
		emojis[name] = parts[1:]
	}

	return emojis
}

var emojiMeta = `---
Title: Emoji
Summary: Emoji Overview
Author: Thomas
Tags:
    - markdown
    - emoji
---
`

var layoutCSS = `
html {
	height: 100%;
}

body {
	height: 100%;
}
`

var layoutHTML = `
{{define "layout"}}
<!doctype html>
<html lang="en">
    <head>
        <meta charset="utf-8">
        {{if .Title}}<title>{{.Title}}</title>{{end}}
        {{if .Description}}<meta name="description" content="{{.Description}}">{{end}}
        {{if .Author}}<meta name="author" content="{{.Author}}">{{end}}
        <link rel="stylesheet" href="/css/layout.css">
        {{range .Layouts}}
            <link rel="stylesheet" href="{{.}}">
        {{end}}
        {{range .Scripts}}
            <script src="{{.}}"></script>
        {{end}} 
    </head>
    <body>
        {{template "content" .}}
    </body>
</html>
{{end}}
`

var markdownHTML = `
{{template "layout" .}}

{{define "content"}}
    {{.Content}}
{{end}}
`

var editHTML = `
{{template "layout" .}}

{{define "content"}}
	<button type="button" id="save">Save</button>
	<input type="hidden" id="name" value="{{.Name}}">
	<hr>
    <textarea id="content">{{.Content}}</textarea>
{{end}}
`

var editCSS = `
#content {
	width: 100%;
	height: 80%;
}
`

var editJS = `
window.onload = function() {
    const saveButton = document.getElementById("save")
    const nameField = document.getElementById("name")
    const contentArea = document.getElementById("content")

    saveButton.addEventListener("click", function() {
        let data = {
            Name: nameField.value,
            Content: contentArea.textContent,
        };

        fetch("/save", {
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json',
            },
            method: "POST",
            body: JSON.stringify(data)
        }).then((response) => {
            response.text().then(function(data) {
                let result = JSON.parse(data);
                console.log("Result", result);
            });
        }).catch((error) => {
            console.log(error);
        });
    })
}
`
