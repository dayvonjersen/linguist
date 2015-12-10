package main

import (
	"log"

	"github.com/generaltso/linguist"
)

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

func getLangFromFilename(filename string) (language string, shouldIgnore bool) {
	language = linguist.DetectFromFilename(filename)
	if language != "" {
		return language, false
	}

	mimetype, shouldIgnore := linguist.DetectMimeFromFilename(filename)
	if mimetype != "" {
		return mimetype, shouldIgnore
	}

	return "", false
}

func getLangFromContents(contents []byte) (language string, shouldIgnore bool) {
	mimetype, shouldIgnore := linguist.DetectMimeFromContents(contents)
	log.Println(mimetype)
	if mimetype != "" {
		return mimetype, shouldIgnore
	}

	language = linguist.DetectFromContents(contents)
	if language != "" {
		return language, false
	}

	return "", false
}
