package config

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/irgangla/markdown-wiki/log"
)

// Config data
type Config struct {
	Port       int
	CommitName string
	CommitMail string
	CommitPush bool
}

// Load config data
func Load() Config {
	var config Config

	var defaults = Config{
		Port:       81,
		CommitName: "thomas",
		CommitMail: "thomas@irgang.eu",
		CommitPush: false,
	}

	data, err := ioutil.ReadFile(filepath.Join(".", "data", "config.json"))
	if err != nil {
		log.Error("CONFIG", "read config file", err)
		return defaults
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Error("CONFIG", "load config data", err)
		return defaults
	}
	return config
}
