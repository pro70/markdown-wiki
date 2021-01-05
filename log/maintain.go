package log

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// DirectoryCheck checks if log dir exists and create it if necessary
func DirectoryCheck() {
	dateTimeFormat := "2006-01-02 15:04:05.000"
	var dir string
	switch runtime.GOOS {
	case "windows":
		{
			executable, err := os.Executable()
			if err != nil {
				fmt.Printf("%v -- ERR -- [LOG] Unable to read actual directory: %v\n", time.Now().Format(dateTimeFormat), err.Error())
			}
			dir = filepath.Dir(executable)
		}
	default:
		{
			executable, err := os.Getwd()
			if err != nil {
				fmt.Printf("%v -- ERR -- [LOG] Unable to read actual directory: %v\n", time.Now().Format(dateTimeFormat), err.Error())
			}
			dir = executable
		}
	}
	logDirectory := filepath.Join(dir, "data", "log")
	_, checkPathError := os.Stat(logDirectory)
	logDirectoryExists := checkPathError == nil
	if logDirectoryExists {
		return
	}
	switch runtime.GOOS {
	case "windows":
		{
			err := os.Mkdir(logDirectory, 0777)
			if err != nil {
				fmt.Printf("%v -- ERR -- [LOG] Unable to create directory for log file: %v\n", time.Now().Format(dateTimeFormat), err.Error())
				return
			}
			fmt.Printf("%v -- INF -- [LOG] Log directory created\n", time.Now().Format(dateTimeFormat))
		}
	default:
		{
			err := os.MkdirAll(logDirectory, 0777)
			if err != nil {
				fmt.Printf("%v -- ERR -- [LOG] Unable to create directory for log file: %v\n", time.Now().Format(dateTimeFormat), err.Error())
				return
			}
			fmt.Printf("%v -- INF -- [LOG] Log directory created\n", time.Now().Format(dateTimeFormat))
		}
	}
}

// DeleteOldLogFiles cleans log folder by removing files which wasn't updated recently
func DeleteOldLogFiles(deleteLogsAfter time.Duration, serviceIsRunning *bool) {
	for *serviceIsRunning {
		directory, err := ioutil.ReadDir("log")
		if err != nil {
			Error("LOG", "Problem opening log directory")
			return
		}
		now := time.Now()
		logDirectory := filepath.Join(".", "log")
		for _, file := range directory {
			if fileAge := now.Sub(file.ModTime()); fileAge > deleteLogsAfter {
				Info("LOG", "Deleting old log file "+file.Name()+" with age of "+fileAge.String())
				logFullPath := filepath.Join(logDirectory, file.Name())
				var err = os.Remove(logFullPath)
				if err != nil {
					Error("LOG", "Problem deleting file "+file.Name()+", "+err.Error())
				}
			}
		}
		time.Sleep(1 * time.Hour)
	}
}
