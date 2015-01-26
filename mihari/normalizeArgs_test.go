package main

import "testing"

var normalizeArgs_original = [][]string{
	[]string{"echo", "hello"},
	[]string{"echo", "hello_hello"},
	[]string{"echo", "hello/hello"},
	[]string{"echo", "hello:hello"},
	[]string{"echo", "hello;hello"},
	[]string{"echo", "hello hello"},
	[]string{"echo", "hello\"hello"},
	[]string{"echo", "hello'hello"},
}

var normalizeArgs_normalized = []string{
	"echo__hello__",
	"echo__hello_Uhello__",
	"echo__hello_Shello__",
	"echo__hello_Chello__",
	"echo__hello_Mhello__",
	"echo__hello_Phello__",
	"echo__hello_Dhello__",
	"echo__hello_Qhello__",
}


func TestNormalizeArgs(t *testing.T) {
	for i, oneArgs := range normalizeArgs_original {
		testValue := NormalizeArgs(oneArgs)
		if testValue != normalizeArgs_normalized[i] {
			t.Errorf("expected: %s  actual: %s", normalizeArgs_normalized[i], testValue)
		}
	}
}

func TestIsEqualArray(t *testing.T) {
	if isEqualArray([]string{"a", "b"}, []string{"a", "b"}) != true {t.Error("Failed 1")}
	if isEqualArray([]string{"a", "b"}, []string{"a", "b", "c"}) != false {t.Error("Failed 2")}
	if isEqualArray([]string{"a", "b"}, []string{"a", "c"}) != false {t.Error("Failed 3")}
}

func TestReverseNormalizeArgs(t *testing.T) {
	for i, oneArgs := range normalizeArgs_normalized {
		testValue := ReverseNormalizedArgs(oneArgs)
		if isEqualArray(testValue, normalizeArgs_original[i]) != true {
			t.Errorf("expected: %s  actual: %s", normalizeArgs_original[i], testValue)
		}
	}
}
