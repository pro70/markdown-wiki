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

	createEmojiList()

	infos, err := ioutil.ReadDir("dummyData")
	if err != nil {
		log.Error("INIT", "dummy data", err)
	}
	for _, i := range infos {
		sourceFile := filepath.Join("dummyData", i.Name())
		destinationFile := filepath.Join("data", "md", i.Name())
		if _, err := os.Stat(destinationFile); os.IsNotExist(err) {
			log.Info("INIT", "copy dummy file", sourceFile)
			input, err := ioutil.ReadFile(sourceFile)
			if err != nil {
				log.Error("INIT", "read source", sourceFile, err)
				continue
			}
			err = ioutil.WriteFile(destinationFile, input, 0775)
			if err != nil {
				log.Error("INIT", "write file", destinationFile, err)
				continue
			}
		}
	}
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

	data, err := ioutil.ReadFile("emoji")
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

var layoutCSS = ``

var layoutHTML = `
{{define "layout"}}
<!doctype html>
<html lang="en">
    <head>
        <meta charset="utf-8">
        {{if .Title}}<title>{{.Title}}</title>{{end}}
        {{if .Description}}<meta name="description" content="{{.Description}}">{{end}}
        {{if .Author}}<meta name="author" content="{{.Author}}">{{end}}
        <link rel="stylesheet" href="css/layout.css">
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
