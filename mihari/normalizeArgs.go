package main

import "bytes"
import "strings"

var normalizeArgs_replaceList = [][]string{
	[]string{"_", "_U"},
	[]string{"/", "_S"},
	[]string{":", "_C"},
	[]string{";", "_M"},
	[]string{" ", "_P"},
	[]string{"\"", "_D"},
	[]string{"'", "_Q"},
}

func NormalizeArgs(args []string) string{
	buffer := bytes.NewBufferString("")

	for _, one := range args {
		for _, oneRep := range normalizeArgs_replaceList {
			one = strings.Replace(one, oneRep[0], oneRep[1], -1)
		}
		buffer.WriteString(one)
		buffer.WriteString("__")
	}

	return buffer.String()
}

func ReverseNormalizedArgs(normalizedArgs string) []string {
	args := strings.Split(normalizedArgs, "__")
	reversedArgs := make([]string, len(args))

	for i, one := range args {
		for _, oneRep := range normalizeArgs_replaceList {
			one = strings.Replace(one, oneRep[1], oneRep[0], -1)
		}
		reversedArgs[i] = one
	}
	return reversedArgs[:len(reversedArgs)-1]
}

func isEqualArray(a []string, b []string) bool {
	if (len(a) != len(b)) {return false}
	for i, one := range a {
		if (one != b[i]) {return false}
	}
	return true
}
