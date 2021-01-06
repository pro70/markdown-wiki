package log

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"time"
)

var (
	writingSync sync.Mutex
	debug       bool
)

// EnableDebug logging
func EnableDebug() {
	debug = true
}

// Info adds an info level message to the log
func Info(reference string, data ...interface{}) {
	message := toMessage(data)
	fmt.Printf("%v -- INF -- [%v] %v\n", time.Now().Format("2006-01-02 15:04:05.000"), reference, message)
	appendDataToLog("INF", reference, message)
}

// Error adds an error level message to the log
func Error(reference string, data ...interface{}) {
	message := toMessage(data)
	fmt.Printf("%v -- ERR -- [%v] %v\n", time.Now().Format("2006-01-02 15:04:05.000"), reference, message)
	appendDataToLog("ERR", reference, message)
}

// Debug adds an debug level message to the log
func Debug(reference string, data ...interface{}) {
	if !debug {
		return
	}
	message := toMessage(data)
	fmt.Printf("%v -- DBG -- [%v] %v\n", time.Now().Format("2006-01-02 15:04:05.000"), reference, message)
	appendDataToLog("DBG", reference, message)
}

func appendDataToLog(logLevel string, reference string, data string) {
	dateTimeFormat := "2006-01-02 15:04:05.000"
	logNameDateTimeFormat := "2006-01-02"
	logDirectory := filepath.Join(".", "data", "log")
	logFileName := fmt.Sprintf("%v.log", time.Now().Format(logNameDateTimeFormat))
	logFullPath := filepath.Join(logDirectory, logFileName)
	logData := fmt.Sprintf("%v %v %v %v", time.Now().Format(dateTimeFormat), reference, logLevel, data)
	writingSync.Lock()
	f, err := os.OpenFile(logFullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("%v -- ERR -- [%v] Cannot open file: %v\n", time.Now().Format(dateTimeFormat), reference, err.Error())
		writingSync.Unlock()
		return
	}
	_, err = f.WriteString(logData + "\n")
	if err != nil {
		fmt.Printf("%v -- ERR -- [%v] Cannot write to file: %v\n", time.Now().Format(dateTimeFormat), reference, err.Error())
	}
	err = f.Close()
	if err != nil {
		fmt.Printf("%v -- ERR -- [%v] Cannot close file: %v\n", time.Now().Format(dateTimeFormat), reference, err.Error())
	}
	writingSync.Unlock()
}

func toMessage(data []interface{}) string {
	message := ""
	for _, d := range data {
		part, ok := d.(string)
		if !ok {
			part = fmt.Sprintf("%v (%v)", d, reflect.TypeOf(d))
		}
		message = message + part + " "
	}
	return message
}
