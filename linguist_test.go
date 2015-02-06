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

func TestDetectFromFilename(t *testing.T) {
	testCases := []struct {
		expected string
		filename string
	}{
		{"C", "foo.c"},
		{"CoffeeScript", "foo.coffee"},
		{"CoffeeScript", "foo._coffee"},
		{"CoffeeScript", "foo.cson"},
		{"CoffeeScript", "foo.iced"},
		{"CoffeeScript", "Cakefile"},
		{"C++", "foo.cpp"},
		{"C++", "foo.C"},
		{"C++", "foo.c++"},
		{"C++", "foo.cxx"},
		{"C++", "foo.H"},
		{"C++", "foo.h++"},
		{"C++", "foo.hpp"},
		{"C++", "foo.hxx"},
		{"C++", "foo.inl"},
		{"C++", "foo.tcc"},
		{"C++", "foo.tpp"},
		{"C++", "foo.ipp"},
		{"CSS", "foo.css"},
		{"HTML", "foo.html"},
		{"HTML", "foo.htm"},
		{"HTML", "foo.xhtml"},
		{"Go", "foo.go"},
		{"Java", "foo.java"},
		{"JavaScript", "foo.js"},
		{"JavaScript", "foo._js"},
		{"JavaScript", "foo.bones"},
		{"JavaScript", "foo.es6"},
		{"JavaScript", "foo.jake"},
		{"JavaScript", "foo.jsfl"},
		{"JavaScript", "foo.jsm"},
		{"JavaScript", "foo.jss"},
		{"JavaScript", "foo.jsx"},
		{"JavaScript", "foo.njs"},
		{"JavaScript", "foo.pac"},
		{"JavaScript", "foo.sjs"},
		{"JavaScript", "foo.ssjs"},
		{"JavaScript", "foo.xsjs"},
		{"JavaScript", "foo.xsjslib"},
		{"JavaScript", "Jakefile"},
		{"Objective-C", "foo.m"},
		{"Objective-C++", "foo.mm"},
		{"Protocol Buffer", "foo.proto"},
		{"Ruby", "foo.rb"},
		{"Ruby", "foo.builder"},
		{"Ruby", "foo.gemspec"},
		{"Ruby", "foo.god"},
		{"Ruby", "foo.irbrc"},
		{"Ruby", "foo.mspec"},
		{"Ruby", "foo.podspec"},
		{"Ruby", "foo.rbuild"},
		{"Ruby", "foo.rbw"},
		{"Ruby", "foo.rbx"},
		{"Ruby", "foo.ru"},
		{"Ruby", "foo.thor"},
		{"Ruby", "foo.watchr"},
		{"Ruby", ".pryrc"},
		{"Ruby", "Appraisals"},
		{"Ruby", "Berksfile"},
		{"Ruby", "Buildfile"},
		{"Ruby", "Gemfile"},
		{"Ruby", "Gemfile.lock"},
		{"Ruby", "Guardfile"},
		{"Ruby", "Jarfile"},
		{"Ruby", "Mavenfile"},
		{"Ruby", "Podfile"},
		{"Ruby", "Puppetfile"},
		{"Ruby", "Thorfile"},
		{"Ruby", "Vagrantfile"},
		{"Ruby", "buildfile"},
		{"Shell", "foo.sh"},
		{"Shell", "foo.bats"},
		{"Shell", "foo.tmux"},
		{"Shell", "Dockerfile"},
		{"SQL", "foo.sql"},
		{"SQL", "foo.prc"},
		{"SQL", "foo.tab"},
		{"SQL", "foo.udf"},
		{"SQL", "foo.viw"},
		{"XML", "foo.xml"},
		{"XML", "foo.rss"},
		{"XML", "foo.svg"},
		{"YAML", "foo.yml"},
		{"YAML", "foo.reek"},
		{"YAML", "foo.rviz"},
		{"YAML", "foo.yaml"},
	}
	for _, c := range testCases {
		l := DetectFromFilename(c.filename)
		if c.expected != l {
			t.Errorf("Expected '%s', but got '%s': %s", c.expected, l, c.filename)
		}
	}
}

func TestDetectFromContents(t *testing.T) {
	testCases := []struct {
		expected string
		contents string
	}{
		{"JavaScript", "#!node"},
		{"Ruby", "#!ruby"},
		{"Shell", "#!bash"},
		{"Shell", "#!sh"},
		{"Shell", "#!zsh"},
	}
	for _, c := range testCases {
		l := DetectFromContents([]byte(c.contents))
		if c.expected != l {
			t.Errorf("Expected '%s', but got '%s': %s", c.expected, l, c.contents)
		}
	}
}

func TestDetectInterpreter(t *testing.T) {
	testCases := []struct {
		expected string
		contents string
	}{
		{"sh", "#!sh"},
		{"sh", "#! sh"},
		{"sh", "#!/bin/sh"},
		{"sh", "#! /bin/sh"},
		{"sh", "#!env sh"},
		{"sh", "#!/bin/env sh"},
		{"sh", "#!sh1"},
		{"sh", "#!sh1.2"},
		{"sh", "#!sh1.2.3"},
		{"env", "#!env"},
	}
	for _, c := range testCases {
		r := detectInterpreter([]byte(c.contents))
		if c.expected != r {
			t.Errorf("Expected '%s', but got '%s': %s", c.expected, r, c.contents)
		}
	}
}
