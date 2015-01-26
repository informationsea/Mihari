package main

import "testing"
import "bytes"
import "time"

func createTestNewRunLogFromExecuteLog() RunLog {
	buf := bytes.NewBufferString("OB/home/hoge:foo.txt")
	buf.WriteByte(0)
	buf.WriteString("OW/home/hoge:foo2.txt")
	buf.WriteByte(0)
	buf.WriteString("OR/home/hoge:foo3.txt")

	runLog, _ := NewRunLogFromExecuteLog(
		[]string{"echo", "ok"},
		time.Date(2015, time.January, 25, 19, 22, 49, 0, time.UTC),
		bytes.NewBuffer(buf.Bytes()), "/home")

	return runLog
}

func TestNewRunLogFromExecuteLog(t *testing.T) {
	runLog := createTestNewRunLogFromExecuteLog()
	if isEqualArray(runLog.Command, []string{"echo", "ok"}) != true {
		t.Error("command is not equal to expected")
	}
	
}

func TestWriteRunLog(t *testing.T) {
	runLog := createTestNewRunLogFromExecuteLog()
	buf := bytes.NewBuffer([]byte{})
	err := WriteRunLog(buf, runLog)
	if err != nil {t.Errorf("Cannot write: %s", err)}
}


func TestReadRunLog(t *testing.T) {
	runLog := createTestNewRunLogFromExecuteLog()
	buf := bytes.NewBuffer([]byte{})
	err := WriteRunLog(buf, runLog)
	if err != nil {t.Errorf("Cannot write: %s", err)}
	newRunLog, err := LoadRunLog(bytes.NewBuffer(buf.Bytes()))
	if isEqualArray(runLog.Command, newRunLog.Command) != true {t.Error("command is not equal to expected")}
	if isEqualArray(runLog.ReadFiles, newRunLog.ReadFiles) != true {t.Error("ReadFiles is not equal to expected")}
	if isEqualArray(runLog.WriteFiles, newRunLog.WriteFiles) != true {t.Error("WriteFiles is not equal to expected")}
	if runLog.LogDate != newRunLog.LogDate {t.Error("LogDate is not equal to expected")}
}

func TestRunLogNormalizePath(t *testing.T) {
	logDir := "/home/foo"
	testValues := []string{"/home/foo:a", "/home/bar:b", "/home/foo:/dev/tty", "/home/bar:/dev/tty", "/home/foo:/home/foo/k"}
	expectedValues := []string{"a", "/home/bar/b", "/dev/tty", "/dev/tty", "k"}

	for i, test := range testValues {
		normalized := runLogNormalizePath(test, logDir)
		if normalized != expectedValues[i] {
			t.Errorf("Expected : %s   Actual : %s", expectedValues[i], normalized)
		}
	}
}
