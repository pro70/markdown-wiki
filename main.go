package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/kardianos/service"

	"github.com/irgangla/markdown-wiki/config"
	"github.com/irgangla/markdown-wiki/events"
	"github.com/irgangla/markdown-wiki/log"
	"github.com/irgangla/markdown-wiki/version"
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
	events.Stop()
	version.Stop()

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
	log.EnableDebug()
	go log.DeleteOldLogFiles(96*time.Hour, &serviceIsRunning)

	initializeData()

	c := config.Load()

	version.Start(c.CommitName, c.CommitMail)

	router := httprouter.New()
	registerRoutes(router)
	events.Start(router)

	host := fmt.Sprintf(":%v", c.Port)
	err := http.ListenAndServe(host, router)
	if err != nil {
		log.Error("RUN", "Problem starting web server:", err.Error())
		os.Exit(-1)
	}
}

func main() {
	log.DirectoryCheck()

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
