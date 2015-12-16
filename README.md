# linguist

Port of [github linguist](https://github.com/github/linguist) to Go. Not complete *but we're getting there...*

Many thanks to [@petermattis](https://github.com/petermattis) for this comment:

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

## Known Issues

1. Using the bayesian classifier should be treated as a last ditch effort, unless you have massive amounts of code on which to train the classifier. Right now it often gives erroneous results and seems to have an unhealthy affinity for the [M language](https://enwp.org/M_programming_language_(disambiguation))

 - more/better programming language sample data

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

#### func  AnalyseWithHints

```go
func AnalyseWithHints(contents []byte, hints []string) (language string)
```

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

#### func  GetColor

```go
func GetColor(language string) string
```
Convenience function that returns the color associated with the language, in
HTML Hex notation (e.g. "#123ABC") from the languages.yaml provided by
github.com/github/linguist

returns empty string if there is no associated color for the language

#### func  GetHintsFromFilename

```go
func GetHintsFromFilename(filename string) (hints []string)
```

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

#### func  IsVendored

```go
func IsVendored(path string) bool
```
IsVendored returns true if path is considered "vendored" and should be excluded
from statistics.

See also the data/vendor.yaml file distributed with this package.

This function also returns true if path is considered documentation.

See also the data/documentation.yaml file distributed with this package.
