package linguist

import (
	"log"
	"os"

	"github.com/generaltso/linguist/tokenizer"
	"github.com/jbrukh/bayesian"
)

var classifier *bayesian.Classifier
var classifier_initialized bool = false

func GetClassifier() *bayesian.Classifier {
	if !classifier_initialized {
		gopath := os.Getenv("GOPATH")
		if gopath == "" {
			gopath = os.Getenv("GOROOT")
			if gopath == "" {
				log.Fatalln("Could not determine $GOPATH or $GOROOT at runtime\n(Necessary to read data dump for linguist)")
			}
		}
		var err error
		classifier, err = bayesian.NewClassifierFromFile(gopath + "/src/github.com/generaltso/linguist/generate_classifier/out/classifier")
		if err != nil {
			log.Panicln(err)
		}
		classifier_initialized = true
	}
	return classifier
}

func Analyse(b []byte) (language string) {
	document := tokenizer.Tokenize(b)
	classifier := GetClassifier()
	_, id, _ := classifier.LogScores(document)

    return string(classifier.Classes[id])
}
