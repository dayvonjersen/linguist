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

var vendorRE *regexp.Regexp

func init() {
	var regexps []string
	bytes := []byte(files["vendor.yaml"])
	if err := yaml.Unmarshal(bytes, &regexps); err != nil {
		log.Fatal(err)
		return
	}

	vendorRE = regexp.MustCompile(strings.Join(regexps, "|"))
}

// IsVendored returns true if path is considered "vendored" and should
// be excluded from statistics.
//
// See also the data/vendor.yaml file distributed with this package.
func IsVendored(path string) bool {
	return vendorRE.MatchString(path)
}
