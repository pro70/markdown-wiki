package main

import (
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/kardianos/service"

	"github.com/irgangla/markdown-wiki/log"
)

const serviceName = "Medium service"
const serviceDescription = "Simple service, just for fun"

var (
	serviceIsRunning bool
	programIsRunning bool
	writingSync      sync.Mutex
)

type program struct{}

func (p program) Start(s service.Service) error {
	log.Info("SERVICE", s.String(), "started")

	writingSync.Lock()
	serviceIsRunning = true
	writingSync.Unlock()

	go p.run()
	return nil
}

func (p program) Stop(s service.Service) error {
	stopSSE()

	writingSync.Lock()
	serviceIsRunning = false
	writingSync.Unlock()

	for programIsRunning {
		log.Info("SERVICE", s.String(), "stopping...")
		time.Sleep(1 * time.Second)
	}

	log.Info("SERVICE", s.String(), "stopped")
	return nil
}

func (p program) run() {
	initializeData()

	log.EnableDebug()
	log.DirectoryCheck()
	go log.DeleteOldLogFiles(96*time.Hour, &serviceIsRunning)

	router := httprouter.New()
	registerRoutes(router)
	startSSE(router)

	err := http.ListenAndServe(":81", router)
	if err != nil {
		log.Error("RUN", "Problem starting web server:", err.Error())
		os.Exit(-1)
	}
}

func main() {
	serviceConfig := &service.Config{
		Name:        serviceName,
		DisplayName: serviceName,
		Description: serviceDescription,
	}
	prg := &program{}
	s, err := service.New(prg, serviceConfig)
	if err != nil {
		log.Error("MAIN", "Cannot create the service: ", err.Error())
	}
	err = s.Run()
	if err != nil {
		log.Error("MAIN", "Cannot start the service: ", err.Error())
	}
}
