package main

import (
	"log"
	"os"
	"path"
)

type LogDirNotFound struct {
	message string
}

func (l LogDirNotFound) Error() string {
	return l.message
}

func searchLogdir() (logdir string, err error) {
	err = nil

	currentPath, lerr := os.Getwd()
	if lerr != nil {log.Fatal(lerr)}
	for currentPath != "/" {
		logCandidate := path.Join(currentPath, ".mihari")
		currentDir, lerr := os.Open(logCandidate)
		if lerr != nil {
			currentPath = path.Clean(path.Join(currentPath, "..", ".."))
			continue
		}

		stat, lerr := currentDir.Stat()
		if lerr != nil {log.Fatal(lerr)}
		if stat.IsDir() {
			logdir = logCandidate
			return
		}
	}
		
	logdir = ""
	err = LogDirNotFound{"Not Found"}
	return
}
