package linguist

import (
	"mime"
	"path/filepath"
	"strings"

	"camlistore.org/pkg/magic"
)

// full mimetype strings to ignore
//
// these are incomplete lists and should be added to
var ignore_mimetype []string = []string{
	"application/octet-stream",
}

// categories of mimetype (the part before the first /) to ignore
//
// these are incomplete lists and should be added to
var ignore_mimetype_start []string = []string{
	"image",
	"audio",
	"video",
}

func shouldIgnoreMime(mimetype string) bool {
	if mimetype == "" {
		return true
	}
	for _, ign := range ignore_mimetype {
		if mimetype == ign {
			return true
		}
	}
	m := strings.Split(mimetype, "/")
	mimetype = m[0]
	for _, ign := range ignore_mimetype_start {
		if mimetype == ign {
			return true
		}
	}
	return false
}

// DetectMimeFromFilename detects the mimetype of the file given by filename
//
// returning the mimetype string, or the empty string on failure
//
// and whether it should be ignored, true if:
//  - true if it is a known binary mimetype and therefore should
//    not be processed by DetectFromContents
//  - true if the mimetype could not be detected
//  - false otherwise
//
// this function will attempt to read the file given by filename
// and will return "", true, and the error if one is encountered,
// otherwise err will be nil
func DetectMimeFromFilename(filename string) (mimetype string, shouldIgnore bool) {
	ext := filepath.Ext(filename)
	if ext != "" {
		by_ext := mime.TypeByExtension(ext)
		if by_ext != "" {
			return by_ext, shouldIgnoreMime(by_ext)
		}
	}
	return "", false
}

func DetectMimeFromContents(contents []byte) (mimetype string, shouldIgnore bool) {
	by_contents := magic.MIMEType(contents)
	if by_contents != "" {
		return by_contents, shouldIgnoreMime(by_contents)
	}
	return "", false
}
