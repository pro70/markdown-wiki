package main

import (
	"net/http"
	"path/filepath"

	"github.com/julienschmidt/httprouter"

	"github.com/irgangla/markdown-wiki/log"
	"github.com/irgangla/markdown-wiki/pages"
	"github.com/irgangla/markdown-wiki/render"
)

func registerRoutes(router *httprouter.Router) {
	router.ServeFiles("/js/*filepath", http.Dir(filepath.Join("data", "js")))
	router.ServeFiles("/css/*filepath", http.Dir(filepath.Join("data", "css")))

	router.GET("/view/:name", render.Endpoint)
	router.GET("/edit/:name", pages.EditEndpoint)

	router.POST("/save", pages.SaveEndpoint)
	router.POST("/log", log.Endpoint)
}
