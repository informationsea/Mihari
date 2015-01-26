package main

import "encoding/json"
import "io"
import "os"
import "bufio"
import "time"
import "strings"
import "path"

type RunLog struct {
	Command []string
	LogDate string
	CurrentDirectory string
	ReadFiles []string
	WriteFiles []string
}

func LoadRunLog(reader io.Reader) (r RunLog, err error) {
	dec := json.NewDecoder(reader)
	err = dec.Decode(&r)
	return
}

func WriteRunLog(writer io.Writer, runLog RunLog) (error) {
	enc := json.NewEncoder(writer)
	err := enc.Encode(&runLog)
	return err
}

func runLogNormalizePath(pathLog string, baseDir string) string {
	components := strings.SplitN(pathLog, ":", 2)
	absolutePath := components[1]
	baseDir = path.Clean(baseDir)
	
	if components[1][0] != '/' {
		absolutePath = path.Clean(path.Join(components[0], components[1]))
	}

	normalizedPath := absolutePath
	if len(absolutePath) > len(baseDir) && absolutePath[:len(baseDir)] == baseDir {
		normalizedPath = absolutePath[len(baseDir)+1:]
	}

	return normalizedPath
}

func NewRunLogFromExecuteLog(args []string, runDate time.Time, log io.Reader, baseDir string) (r RunLog, err error) {
	scanner := bufio.NewScanner(log)
	scanner.Split(ScanNull)
	currentDir, _ := os.Getwd()
	r = RunLog{Command:args, LogDate:runDate.String(), CurrentDirectory:currentDir}
	
	for scanner.Scan() {
		one := scanner.Text()
		if one[0] != 'O' {continue}
		switch (one[1]) {
		case 'B':
			r.ReadFiles = append(r.ReadFiles, runLogNormalizePath(one[2:], baseDir))
			r.WriteFiles = append(r.WriteFiles, runLogNormalizePath(one[2:], baseDir))
		case 'W':
			r.WriteFiles = append(r.WriteFiles, runLogNormalizePath(one[2:], baseDir))
		case 'R':
			r.ReadFiles = append(r.ReadFiles, runLogNormalizePath(one[2:], baseDir))
		}
	}
	return
}
