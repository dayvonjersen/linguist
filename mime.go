package linguist

import (
	"mime"
	"path/filepath"
	"strings"

	"github.com/rakyll/magicmime"
)

// full mimetype strings to ignore
//
// NOTE(tso): these are incomplete lists and should be added to
var ignore_mimetype []string = []string{
	"application/octet-stream",
}

// categories of mimetype (the part before the first /) to ignore
//
// NOTE(tso): these are incomplete lists and should be added to
var ignore_mimetype_start []string = []string{
	"image",
	"audio",
	"video",
}

func shouldIgnoreMime(mimetype string) bool {
	if mimetype == "" {
		return false
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
// Returns the mimetype string, or the empty string on failure
//
// shouldIgnore will be true iff the mimetype matches known binary formats
//
// This function uses the golang.org/pkg/mime library
// and should be relatively safe to use, but not very robust
func DetectMimeFromFilename(filename string) (mimetype string, shouldIgnore bool) {
	ext := filepath.Ext(filename)
	if ext != "" {
		by_ext := mime.TypeByExtension(ext)
		if by_ext != "" {
			splix := strings.Split(by_ext, ";")
			mr_mime := strings.TrimSpace(splix[0])
			return mr_mime, shouldIgnoreMime(mr_mime)
		}
	}
	return "", false
}

// DetectMimeFromContents detects the mimetype based on the contents given
//
// Returns the mimetype string, or the empty string on failure
//
// shouldIgnore will be true iff the mimetype matches known binary formats
//
// This function uses the github.com/rakyll/magicmime library
// and may not be compatible with your system
func DetectMimeFromContents(contents []byte) (mimetype string, shouldIgnore bool) {
	magicmime.Open(magicmime.MAGIC_MIME_TYPE)
	defer magicmime.Close()
	by_contents, err := magicmime.TypeByBuffer(contents)
	if err != nil {
		println(err.Error())
	}
	if by_contents != "" {
		return by_contents, shouldIgnoreMime(by_contents)
	}
	return "", false
}
