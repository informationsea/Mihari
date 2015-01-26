package main

import (
	"testing"
	"bufio"
	"strings"
	"bytes"
)

func TestScanNull1(t *testing.T) {
	testText := "aabb"
	scanner := bufio.NewScanner(strings.NewReader(testText))
	scanner.Split(ScanNull)

	if scanner.Scan() != true {t.Error("First end?")}

	testValue := scanner.Text()
	if testValue != "aabb" {t.Errorf("Not expected value: %s", testValue)}
	
	if scanner.Scan() != false {t.Error("not end?")}
}

func TestScanNull2(t *testing.T) {
	testText := bytes.NewBufferString("aabb")
	testText.WriteByte(0x00)
	testText.WriteString("cc")
	testText.WriteByte(0x00)
	testText.WriteString("dd")

	scanner := bufio.NewScanner(bytes.NewBuffer(testText.Bytes()))
	scanner.Split(ScanNull)

	expectedValues := []string{"aabb", "cc", "dd"}
	count := 0
	for scanner.Scan() {
		if (len(expectedValues) <= count) {
			t.Errorf("out of range %d", count)
			t.FailNow()
		}
		
		
		testValue := scanner.Text()
		if testValue != expectedValues[count] {
			t.Errorf("expected: %s / actural: %s [%d]", expectedValues[count],  testValue, len(testValue))
		}
		count += 1
	}

	if count != 3 {t.Error("Count is not valid")}
}
