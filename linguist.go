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

	bytes := []byte(Files["languages.yaml"])
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

// DetectFromFilename detects the language solely from the filename,
// returning the empty string on ambiguous or unknown filenames.
func DetectFromFilename(filename string) string {
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

// DetectFromContents detects the language from the file contents,
// returning the empty string if the language could not be determined.
func DetectFromContents(contents []byte) string {
	interpreter := detectInterpreter(contents)
	if interpreter != "" {
		if l := interpreters[interpreter]; len(l) == 1 {
			return l[0]
		}
	}
	// TODO(pmattis): Linguist falls back to using a bayesian classifier
	// at this point. Wouldn't be hard too do something similar using
	// their classification data (which is stored in the samples.json
	// file). Need to do this to properly detect the language for .h
	// files (C, C++, Objective-C, Objective-C++).
	return Analyse(contents)
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
