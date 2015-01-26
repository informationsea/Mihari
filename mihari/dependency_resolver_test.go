package main

import ("testing"
	"time"
	"fmt"
	"bytes")

var __depdency_resolver_test_dates = []string{time.Date(2015, time.January, 25, 23, 24, 10, 23, time.UTC).String(),}
var __depdency_resolver_test_runlogs = []RunLog{
	RunLog{Command: []string{"a", "b"}, LogDate: __depdency_resolver_test_dates[0],
		ReadFiles:[]string{"1"}, WriteFiles:[]string{"2"}},
	RunLog{Command: []string{"c", "d"}, LogDate: __depdency_resolver_test_dates[0],
		ReadFiles:[]string{"2"}, WriteFiles:[]string{"3"}},
	RunLog{Command: []string{"e", "f"}, LogDate: __depdency_resolver_test_dates[0],
		ReadFiles:[]string{"1", "3"}, WriteFiles:[]string{"4", "5"}},
	RunLog{Command: []string{"g", "h"}, LogDate: __depdency_resolver_test_dates[0],
		ReadFiles:[]string{"1", "2"}, WriteFiles:[]string{"6"}},
}


func TestResolveDependency(t *testing.T) {
	runlogs := __depdency_resolver_test_runlogs

	files, commands, err := ResolveDependency(runlogs)

	if err != nil {t.Errorf("Resolve error %s", err)}

	if len(files) != 6 {t.Error("# of files is not equal to an expected value")}
	if len(commands) != 4 {t.Error("# of commands is not equal to an expected value")}

	fmt.Println(files)
	fmt.Println(commands)
}


func TestGenerateMakefile(t *testing.T) {
	files, commands, err := ResolveDependency(__depdency_resolver_test_runlogs)
	if err != nil {t.Errorf("Resolve error %s", err)}

	buf := bytes.NewBuffer([]byte{})
	err = generateMakefile(files, commands, buf)

	fmt.Println(buf.String())
}

func TestEscapeCommand(t *testing.T) {
	testValue := [][]string{
		[]string{"a", "b", "c"},
		[]string{"a", "b c", "d"},
		[]string{"a", "b\"c", "c"},
		[]string{"a", "b\" c", "c"},
	}
	expected := []string{
		"a b c",
		"a \"b c\" d",
		"a b\\\"c c",
		"a \"b\\\" c\" c",
	}

	for i, one := range testValue {
		treated := escapeCommand(one)
		if treated != expected[i] {
			t.Errorf("expected: %s  actual: %s", expected[i], treated)
		}
	}
	
}
