package linguist

import (
	"log"
	"os"

	"github.com/generaltso/linguist/tokenizer"
	"github.com/jbrukh/bayesian"
)

var classifier *bayesian.Classifier
var classifier_initialized bool = false

// Gets the baysian.Classifier which has been trained on programming language
// samples from github.com/github/linguist after running the generator found
// in cmd/generate_classifier/main.go
//
// NOTE(tso): github.com/jbrukh/bayesian provides the mechanism for
// serialization/deserialization of Classifier objects. After many failed
// attempts to serialize this data *directly* into this package, both as a
// base64 encoded string and exporting the structure from bayesian.Classifier
// with an old-fashioned %#v (really bad idea), I have resorted to using the
// serialization to file option, which is not very portable.
//
// Currently it looks for the data dump distributed with this package
// using the GOPATH or GOROOT env vars at runtime.
func getClassifier() *bayesian.Classifier {
	if !classifier_initialized {
		gopath := os.Getenv("GOPATH")
		if gopath == "" {
			gopath = os.Getenv("GOROOT")
			if gopath == "" {
				log.Fatalln("Could not determine GOPATH or GOROOT at runtime\n(Necessary to read data dump for linguist)")
			}
		}
		var err error
		classifier, err = bayesian.NewClassifierFromFile(gopath + "/src/github.com/generaltso/linguist/data/classifier")
		if err != nil {
			log.Panicln(err)
		}
		classifier_initialized = true
	}
	return classifier
}

// Attempts to use Naive Bayesian Classification on the file contents provided
//
// Returns the name of a programming language, or the empty string if one could
// not be determined.
//
// NOTE(tso): May yield inaccurate results
func Analyse(contents []byte) (language string) {
	document := tokenizer.Tokenize(contents)
	classifier := getClassifier()
	_, id, _ := classifier.LogScores(document)
	return string(classifier.Classes[id])
}
