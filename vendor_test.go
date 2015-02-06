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

import "testing"

func TestIsVendored(t *testing.T) {
	testCases := []struct {
		path     string
		expected bool
	}{
		{"cache/bar", true},
		{"foo/cache/bar", true},
		{"deps/bar", true},
		{"foo/deps/bar", false},
		{"tools/bar", true},
		{"foo/tools/bar", false},
		{"tools/bar", true},
		{"bootstrap.js", true},
		{"bootstrap.min.js", true},
		{"jquery.js", true},
		{"foo.pb.go", true},
		{"foo.pb.goo", false},
	}
	for _, c := range testCases {
		v := IsVendored(c.path)
		if c.expected != v {
			t.Errorf("Expected %v, but got %v: %s", c.expected, v, c.path)
		}
	}
}
