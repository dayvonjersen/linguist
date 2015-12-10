package main

import "github.com/generaltso/linguist"

var langs map[string]int = make(map[string]int)
var total_size int = 0
var res map[string]int = make(map[string]int)
var num_files int = 0
var max_len int = 0

func putResult(language string, size int) {
	res[language]++
	langs[language] += size
	total_size += size
	num_files++
	if len(language) > max_len {
		max_len = len(language)
	}
}

func getLangFromFilename(filename string) string {
	res1 := linguist.DetectFromFilename(filename)
	if res1 != "" {
		return res1
	}

	mimetype, shouldIgnore := linguist.DetectMimeFromFilename(filename)
	if shouldIgnore {
		return mimetype
	}

	return ""
}

func getLangFromContents(contents []byte) string {
	//	mimetype, shouldIgnore := linguist.DetectMimeFromContents(contents)
	//	if shouldIgnore {
	//		return mimetype
	//	}

	language := linguist.DetectFromContents(contents)
	if language != "" {
		return language
	}

	return ""
}
