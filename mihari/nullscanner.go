package main

import "bytes"

func ScanNull(data []byte, atEOF bool) (advance int, token []byte, error error){
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if i := bytes.IndexByte(data, 0); i >= 0 {
		advance = i+1
		token = data[:i]
		return
	}
	
	if atEOF {
		advance = len(data)
		token = data
	}
	return
}
