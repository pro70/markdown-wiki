package main

import (
	"encoding/json"

	"github.com/julienschmidt/httprouter"
	"github.com/julienschmidt/sse"

	"github.com/irgangla/markdown-wiki/log"
)

// Event for SSE sending
type Event struct {
	// ID of the event
	ID string
	// Event name
	Event string
	// Data of the event
	Data interface{}
}

var (
	events chan Event
)

func startSSE(router *httprouter.Router) *chan Event {
	events = make(chan Event)
	sender := sse.New()
	router.Handler("GET", "/event", sender)
	go streamEvents(sender)
	return &events
}

func stopSSE() {
	close(events)
}

func streamEvents(sender *sse.Streamer) {
	log.Info("SSE", "Streaming events started")
	for e := range events {
		log.Debug("SSE", "Send event", e)
		sender.SendString(e.ID, e.Event, marshal(e.Data))
	}
}

func marshal(d interface{}) string {
	data, err := json.Marshal(d)
	if err != nil {
		log.Error("SSE", "marshal data", err)
		return ""
	}
	return string(data)
}
