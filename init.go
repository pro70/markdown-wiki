package main

import (
	"os"
	"path/filepath"
)

func initializeData() {
	makeDirs()
}

func makeDirs() {
	dirs := [][]string{
		[]string{".", "data", "css"},
		[]string{".", "data", "html"},
		[]string{".", "data", "js"},
		[]string{".", "data", "md"},
		[]string{".", "data", "template"},
		[]string{".", "data", "template", "shared"},
	}

	for _, d := range dirs {
		path := filepath.Join(d...)
		os.MkdirAll(path, 0775)
	}
}
