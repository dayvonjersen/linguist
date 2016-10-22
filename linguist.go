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
	"bufio"
	"bytes"
	"log"
	"path/filepath"
	"regexp"

	"gopkg.in/yaml.v1"
)

var (
	extensions   = map[string][]string{}
	filenames    = map[string][]string{}
	interpreters = map[string][]string{}
	colors       = map[string]string{}

	shebangRE       = regexp.MustCompile(`^#!\s*(\S+)(?:\s+(\S+))?.*`)
	scriptVersionRE = regexp.MustCompile(`((?:\d+\.?)+)`)
)

func init() {
	type language struct {
		Extensions   []string `yaml:"extensions,omitempty"`
		Filenames    []string `yaml:"filenames,omitempty"`
		Interpreters []string `yaml:"interpreters,omitempty"`
		Color        string   `yaml:"color,omitempty"`
	}
	languages := map[string]*language{}

	bytes := []byte(files["data/languages.yml"])
	if err := yaml.Unmarshal(bytes, languages); err != nil {
		log.Fatal(err)
	}

	for n, l := range languages {
		for _, e := range l.Extensions {
			extensions[e] = append(extensions[e], n)
		}
		for _, f := range l.Filenames {
			filenames[f] = append(filenames[f], n)
		}
		for _, i := range l.Interpreters {
			interpreters[i] = append(interpreters[i], n)
		}
		colors[n] = l.Color
	}
}

// Convenience function that returns the color associated
// with the language, in HTML Hex notation (e.g. "#123ABC")
// from the languages.yml file provided by https://github.com/github/linguist
//
// Returns the empty string if there is no associated color for the language.
func LanguageColor(language string) string {
	if c, ok := colors[language]; ok {
		return c
	}
	return ""
}

// Attempts to determine the language of a source file based solely on 
// common naming conventions and file extensions
// from the languages.yml file provided by https://github.com/github/linguist
//
// Returns the empty string in ambiguous or unrecognized cases.
func LanguageByFilename(filename string) string {
	if l := filenames[filename]; len(l) == 1 {
		return l[0]
	}
	ext := filepath.Ext(filename)
	if ext != "" {
		if l := extensions[ext]; len(l) == 1 {
			return l[0]
		}
	}
	return ""
}

// Attempts to detect all possible languages of a source file based solely on 
// common naming conventions and file extensions
// from the languages.yml file provided by https://github.com/github/linguist
//
// Intended to be used with LanguageByContents.
//
// May return an empty slice.
func LanguageHints(filename string) (hints []string) {
	if l, ok := filenames[filename]; ok {
		hints = append(hints, l...)
	}
	if ext := filepath.Ext(filename); ext != "" {
		if l, ok := extensions[ext]; ok {
			hints = append(hints, l...)
		}
	}
	return hints
}

// Attempts to detect the language of a source file based on its
// contents and a slice of hints to the possible answer.
//
// Obtain hints with LanguageHints()
//
// Returns the empty string a language could not be determined.
func LanguageByContents(contents []byte, hints []string) string {
	interpreter := detectInterpreter(contents)
	if interpreter != "" {
		if l := interpreters[interpreter]; len(l) == 1 {
			return l[0]
		}
	}
	return Analyse(contents, hints)
}

func detectInterpreter(contents []byte) string {
	scanner := bufio.NewScanner(bytes.NewReader(contents))
	scanner.Scan()
	line := scanner.Text()
	m := shebangRE.FindStringSubmatch(line)
	if m == nil || len(m) != 3 {
		return ""
	}
	base := filepath.Base(m[1])
	if base == "env" && m[2] != "" {
		base = m[2]
	}
	// Strip suffixed version number.
	return scriptVersionRE.ReplaceAllString(base, "")
}
