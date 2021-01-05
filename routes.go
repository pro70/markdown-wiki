package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/irgangla/markdown-wiki/log"
	"github.com/irgangla/markdown-wiki/render"
)

func registerRoutes(router *httprouter.Router) {
	router.ServeFiles("/js/*filepath", http.Dir("js"))
	router.ServeFiles("/css/*filepath", http.Dir("css"))

	router.GET("/view/:name", render.Endpoint)

	router.POST("/log", log.Endpoint)
}
