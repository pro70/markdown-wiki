package events

import (
	"encoding/json"

	"github.com/julienschmidt/httprouter"
	"github.com/julienschmidt/sse"

	"github.com/irgangla/markdown-wiki/log"
	"github.com/irgangla/markdown-wiki/sdk"
)

// Start starts server side event sending
func Start(router *httprouter.Router) {
	sdk.ClientEvents = make(chan sdk.Event)
	sender := sse.New()
	router.Handler("GET", "/event", sender)
	go streamEvents(sender)
}

// Stop stops server side event sending
func Stop() {
	close(sdk.ClientEvents)
}

func streamEvents(sender *sse.Streamer) {
	log.Info("SSE", "Streaming events started")
	for e := range sdk.ClientEvents {
		log.Debug("SSE", "Send event", e)
		sender.SendString(e.ID, e.Event, marshal(e.Data))
	}
	log.Info("SSE", "Streaming events stopped")
}

func marshal(d interface{}) string {
	data, err := json.Marshal(d)
	if err != nil {
		log.Error("SSE", "marshal data", err)
		return ""
	}
	return string(data)
}
