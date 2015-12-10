# linguist

Port of [github linguist](https://github.com/github/linguist) to Go. Not complete **but we're getting there...**

Many thanks to @petermattis for this comment:

```
	// TODO(pmattis): Linguist falls back to using a bayesian classifier
	// at this point. Wouldn't be hard too do something similar using
	// their classification data (which is stored in the samples.json
	// file). Need to do this to properly detect the language for .h
	// files (C, C++, Objective-C, Objective-C++).
```

Which allowed me to find this library:

[github.com/jbrukh/bayesian](https://github.com/jbrukh/bayesian)

You can blame that and/or the amateurish tokenizer I wrote for any inaccurate results.

See cmd/l for a reference implementation command-line tool.

###### [godocdown](https://github.com/robertkrimen/godocdown) >> README.md
# linguist
--
    import "github.com/generaltso/linguist"


## Usage

#### func  Analyse

```go
func Analyse(contents []byte) (language string)
```
Attempts to use Naive Bayesian Classification on the file contents provided

Returns the name of a programming language, or the empty string if one could not
be determined.

NOTE(tso): May yield inaccurate results

#### func  DetectFromContents

```go
func DetectFromContents(contents []byte) string
```
DetectFromContents detects the language from the file contents, returning the
empty string if the language could not be determined.

#### func  DetectFromFilename

```go
func DetectFromFilename(filename string) string
```
DetectFromFilename detects the language solely from the filename, returning the
empty string on ambiguous or unknown filenames.

#### func  DetectMimeFromContents

```go
func DetectMimeFromContents(contents []byte) (mimetype string, shouldIgnore bool)
```
DetectMimeFromContents detects the mimetype based on the contents given

Returns the mimetype string, or the empty string on failure

shouldIgnore will be true iff the mimetype matches known binary formats

#### func  DetectMimeFromFilename

```go
func DetectMimeFromFilename(filename string) (mimetype string, shouldIgnore bool)
```
DetectMimeFromFilename detects the mimetype of the file given by filename

Returns the mimetype string, or the empty string on failure

shouldIgnore will be true iff the mimetype matches known binary formats

#### func  GetColor

```go
func GetColor(language string) string
```
Convenience function that returns the color associated with the language, in
HTML Hex notation (e.g. "#123ABC") from the languages.yaml provided by
github.com/github/linguist

returns empty string if there is no associated color for the language

#### func  IsVendored

```go
func IsVendored(path string) bool
```
IsVendored returns true if path is considered "vendored" and should be excluded
from statistics.

See also the data/vendor.yaml file distributed with this package.
