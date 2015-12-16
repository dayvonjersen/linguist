// Copyright 2015 Peter Mattis.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License. See the AUTHORS file
// for names of contributors.

package linguist

import (
	"log"
	"regexp"
	"strings"

	"gopkg.in/yaml.v1"
)

// ShouldIgnoreFilename returns true if filename matches known files which
// typically should not be passed to LanguageByFilename.
//
// Right now, this simply calls IsVendored and IsDocumentation.
func ShouldIgnoreFilename(filename string) bool {
	return IsVendored(filename) || IsDocumentation(filename)
}

// ShouldIgnoreContents returns true if contents match known files which
// typically should not be passed to LangugeByContents.
//
// Right now, this simply calls IsBinary.
func ShouldIgnoreContents(contents []byte) bool {
	return IsBinary(contents)
}

var vendorRE *regexp.Regexp
var doxRE *regexp.Regexp

func init() {
	var regexps []string
	bytes := []byte(files["data/vendor.yaml"])
	if err := yaml.Unmarshal(bytes, &regexps); err != nil {
		log.Fatal(err)
		return
	}
	vendorRE = regexp.MustCompile(strings.Join(regexps, "|"))

	var moreregex []string
	bytes = []byte(files["data/documentation.yaml"])
	if err := yaml.Unmarshal(bytes, &moreregex); err != nil {
		log.Fatal(err)
		return
	}
	doxRE = regexp.MustCompile(strings.Join(moreregex, "|"))
}

// IsVendored returns true if path is considered "vendored" and should
// be excluded from statistics.
//
// See also the data/vendor.yaml file distributed with this package.
func IsVendored(path string) bool {
	return vendorRE.MatchString(path)
}

// IsDocumentation returns true if path is considered documentation.
//
// See also the data/documentation.yaml file distributed with this package.
func IsDocumentation(path string) bool {
	return doxRE.MatchString(path)
}

// IsBinary checks contents for known character escape codes which
// frequently show up in binary files but rarely (if ever) in text.
//
// Use this check before using DetectFromContents to reduce likelihood
// of passing binary data into it.
//
// NOTE(tso): preliminary testing on this method of checking for binary
// contents were promising, having fed a document consisting of all
// utf-8 codepoints from 0000 to FFFF with satisfactory results. Thanks
// to robpike.io/cmd/unicode:
// ```
// unicode -c $(seq 0 65535 | xargs printf "%04x ") | tr -d '\n' > unicode_test
// ```
//
// However, the intentional presence of character escape codes to throw
// this function off is entirely possible, as is, potentially, a binary
// file consisting entirely of the 4 exceptions to the rule for the first
// 512 bytes. It is also possible that more character escape codes need
// to be added.
//
// Further analysis and real world testing of this is required.
func IsBinary(contents []byte) (probably bool) {
	for _, b := range contents[:512] {
		if b < 32 {
			switch b {
			case 0:
				fallthrough
			case 9:
				fallthrough
			case 10:
				fallthrough
			case 13:
				continue
			default:
				return true
			}
		}
	}
	return false
}
