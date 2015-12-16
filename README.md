# linguist

<s>Port</s> *reimagining* of [github linguist](https://github.com/github/linguist) in Go. 

Many thanks to [@petermattis](https://github.com/petermattis) for his initial work in creating this project, and especially thanks this comment:

```
	// TODO(pmattis): Linguist falls back to using a bayesian classifier
	// at this point. Wouldn't be hard too do something similar using
	// their classification data (which is stored in the samples.json
	// file). Need to do this to properly detect the language for .h
	// files (C, C++, Objective-C, Objective-C++).
```

Which allowed me to find this library:

[github.com/jbrukh/bayesian](https://github.com/jbrukh/bayesian)

You can blame that, or more likely the amateurish `tokenizer` I made, for any inaccurate results.

See `cmd/l` for a reference implementation in the form of a command-line tool.

## Known Issues

1. Using the bayesian classifier should be treated as a last ditch effort, unless you have massive amounts of code on which to train the classifier. Right now it often gives erroneous results and seems to have an unhealthy affinity for the [M language](https://enwp.org/M_programming_language_(disambiguation))

 - more/better programming language sample data

###### [godocdown](https://github.com/robertkrimen/godocdown) >> README.md
# linguist
--
    import "github.com/generaltso/linguist"


## Usage

#### func  IsBinary

```go
func IsBinary(contents []byte) (probably bool)
```
IsBinary checks contents for known character escape codes which frequently show
up in binary files but rarely (if ever) in text.

Use this check before using DetectFromContents to reduce likelihood of passing
binary data into it.

NOTE(tso): preliminary testing on this method of checking for binary contents
were promising, having fed a document consisting of all utf-8 codepoints from
0000 to FFFF with satisfactory results. Thanks to robpike.io/cmd/unicode: ```
unicode -c $(seq 0 65535 | xargs printf "%04x ") | tr -d '\n' > unicode_test ```

However, the intentional presence of character escape codes to throw this
function off is entirely possible, as is, potentially, a binary file consisting
entirely of the 4 exceptions to the rule for the first 512 bytes. It is also
possible that more character escape codes need to be added.

Further analysis and real world testing of this is required.

#### func  IsDocumentation

```go
func IsDocumentation(path string) bool
```
IsDocumentation returns true if path is considered documentation.

See also the data/documentation.yaml file distributed with this package.

#### func  IsVendored

```go
func IsVendored(path string) bool
```
IsVendored returns true if path is considered "vendored" and should be excluded
from statistics.

See also the data/vendor.yaml file distributed with this package.

#### func  LanguageByContents

```go
func LanguageByContents(contents []byte, hints []string) string
```
LanguageByContents attempts to detect the language based on its contents and a
slice of hints to the possible answer (obtained with LanguageHints), returning
the empty string a language could not be determined.

#### func  LanguageByFilename

```go
func LanguageByFilename(filename string) string
```
LanguageByFilename attempts to detect the language of a file, based only on its
name, returning the empty string in ambiguous or unrecognized cases.

#### func  LanguageColor

```go
func LanguageColor(language string) string
```
Convenience function that returns the color associated with the language, in
HTML Hex notation (e.g. "#123ABC") from the languages.yaml provided by
github.com/github/linguist

returns empty string if there is no associated color for the language

#### func  LanguageHints

```go
func LanguageHints(filename string) (hints []string)
```
Similarly to LanguageByFilename, LanguageHints attempts to detect the language
of a file based solely on its name, returning all known possiblities as a slice
of strings.

Intended to be used with LanguageByContents.

#### func  ShouldIgnoreContents

```go
func ShouldIgnoreContents(contents []byte) bool
```
ShouldIgnoreContents returns true if contents match known files which typically
should not be passed to LangugeByContents.

Right now, this simply calls IsBinary.

#### func  ShouldIgnoreFilename

```go
func ShouldIgnoreFilename(filename string) bool
```
ShouldIgnoreFilename returns true if filename matches known files which
typically should not be passed to LanguageByFilename.

Right now, this simply calls IsVendored and IsDocumentation.
