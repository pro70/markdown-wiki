package log

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

// ClientLog request json struct
type ClientLog struct {
	Level   string
	Tag     string
	Message string
}

// ClientLogResponse json struct
type ClientLogResponse struct {
	Result   string
	Duration string
}

// Endpoint is the POST endpoint for client log requests
func Endpoint(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	var data ClientLog
	timer := time.Now()

	Info("CLIENT", "Log message from client", request.RemoteAddr)

	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		Error("CLIENT", "request data error", request.RemoteAddr, err)
		var responseData ClientLogResponse
		responseData.Result = "nok"
		writer.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(writer).Encode(responseData)
		if err != nil {
			Error("CLIENT", err.Error())
		}
		return
	}

	switch data.Level {
	case "INF":
		Info(data.Tag, data.Message)
	default:
		Error(data.Tag, data.Message)
	}

	end := time.Since(timer)
	Info("CLIENT", request.RemoteAddr, "processing takes:", end.String())

	var responseData ClientLogResponse
	responseData.Result = "ok"
	responseData.Duration = end.String()

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(responseData)
	if err != nil {
		Error("CLIENT", "response encode error", err.Error())
	}
}
