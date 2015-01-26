package main

import ("time"
	"strings"
	"io"
	"bufio")


type MihariFile struct {
	FileName string
	GeneratorCommand []*MihariCommand
	DependentCommand []*MihariCommand
}

func JoinMihariFile(a []*MihariFile, separator string) (text string) {
	text = ""
	for i, one := range a {
		if i != 0 {text += separator}
		text += one.FileName
	}
	return
}

type MihariCommand struct {
	Command []string
	LastRunDate time.Time
	Log RunLog
	GenerateFile []*MihariFile
	DependedFile []*MihariFile
	ShouldUse bool
}

func _createNewFile(fileName string) *MihariFile {
	return &MihariFile{
		FileName: fileName,
		GeneratorCommand: []*MihariCommand{},
		DependentCommand: []*MihariCommand{},
	}
}

type CannotResolveDependencyError struct {
	message string
}

func (e CannotResolveDependencyError) Error() string {
	return e.message
}

func ResolveDependency(logs []RunLog) (files []MihariFile, commands []MihariCommand, err error) {
	files = []MihariFile{}
	fileMap := make(map[string]*MihariFile)
	commands = []MihariCommand{}
	err = nil

	const timeFormat = "2006-01-02 15:04:05.999999999 -0700 MST"

	for _, oneLog := range logs {
		t, _ := time.Parse(timeFormat, oneLog.LogDate)
		currentCommand := MihariCommand{
			Command: oneLog.Command, Log: oneLog, LastRunDate: t,
			GenerateFile: []*MihariFile{}, DependedFile: []*MihariFile{},
			ShouldUse: true,
		}

		for _, oneFile := range oneLog.ReadFiles {
			if oneFile[0] == '/' {continue} // skip outside of working directory
			_, present := fileMap[oneFile]
			if !present {fileMap[oneFile] = _createNewFile(oneFile)}
			
			fileMap[oneFile].DependentCommand = append(fileMap[oneFile].DependentCommand, &currentCommand)
			currentCommand.DependedFile = append(currentCommand.DependedFile, fileMap[oneFile])
		}

		for _, oneFile := range oneLog.WriteFiles {
			if oneFile[0] == '/' {continue} // skip outside of working directory
			_, present := fileMap[oneFile]
			if !present {fileMap[oneFile] = _createNewFile(oneFile)}
			
			fileMap[oneFile].GeneratorCommand = append(fileMap[oneFile].GeneratorCommand, &currentCommand)
			currentCommand.GenerateFile = append(currentCommand.GenerateFile, fileMap[oneFile])
		}
		commands = append(commands, currentCommand)
	}

	for _, v := range fileMap {
		if len(v.GeneratorCommand) > 1 {
			files = nil
			commands = nil
			err = CannotResolveDependencyError{message: "Failed to resolve dependency"}
			return
		}
	}

	for _, v := range fileMap {
		files = append(files, *v)
	}
	
	return
}

func escapeCommand(command []string) (result string) {
	result = ""
	for i, v := range command {
		if i != 0 {
			result += " "
		}

		v = strings.Replace(v, "\"", "\\\"", -1)
		
		if strings.Contains(v, " ") {
			result += "\"" + v + "\""
		} else {
			result += v
		}
	}
	return
}

func generateMakefile(files []MihariFile, commands []MihariCommand, writer io.Writer) (err error) {

	outputFiles := []string{}
	inputFiles := []string{}
	generatedFiles := []string{}
	err = nil

	for _, oneFile := range files {
		if len(oneFile.DependentCommand) == 0 {
			outputFiles = append(outputFiles, oneFile.FileName)
		}
		if len(oneFile.GeneratorCommand) == 0 {
			inputFiles = append(inputFiles, oneFile.FileName)
		}
		if len(oneFile.GeneratorCommand) > 0 {
			generatedFiles = append(generatedFiles, oneFile.FileName)
		}
	}

	makefileWriter := bufio.NewWriter(writer)
	makefileWriter.WriteString("all: "+strings.Join(outputFiles, " ")+"\n\n")
	makefileWriter.WriteString("clean:\n")
	makefileWriter.WriteString("\t-rm "+strings.Join(generatedFiles, " ")+"\n\n")
	makefileWriter.WriteString(".PHONY: all clean\n\n")

	for _, oneCommand := range commands {
		if !oneCommand.ShouldUse {continue}
		makefileWriter.WriteString(JoinMihariFile(oneCommand.GenerateFile, " ") + ":" + JoinMihariFile(oneCommand.DependedFile, " ")+"\n")
		makefileWriter.WriteString("\t"+escapeCommand(oneCommand.Command)+"\n\n")
	}
	makefileWriter.Flush()

	return
}
