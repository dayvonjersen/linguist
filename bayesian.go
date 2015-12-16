package linguist

import (
	"bytes"
	"log"
	"math"

	"github.com/generaltso/linguist/data"
	"github.com/generaltso/linguist/tokenizer"
	"github.com/jbrukh/bayesian"
)

var classifier *bayesian.Classifier
var classifier_initialized bool = false

// Gets the baysian.Classifier which has been trained on programming language
// samples from github.com/github/linguist after running the generator found
// in cmd/generate_classifier/main.go
func getClassifier() *bayesian.Classifier {
	if !classifier_initialized {
		data, err := data.Asset("classifier")
		if err != nil {
			log.Panicln(err)
		}
		reader := bytes.NewReader(data)
		classifier, err = bayesian.NewClassifierFromReader(reader)
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

func AnalyseWithHints(contents []byte, hints []string) (language string) {
	document := tokenizer.Tokenize(contents)
	classifier := getClassifier()
	scores, _, _ := classifier.LogScores(document)

	langs := map[string]struct{}{}
	for _, hint := range hints {
		langs[hint] = struct{}{}
	}
	best_score := math.Inf(-1)
	best_answer := ""

	for id, score := range scores {
		answer := string(classifier.Classes[id])
		if _, ok := langs[answer]; ok {
			if score >= best_score {
				best_score = score
				best_answer = answer
			}
		}
	}
	return best_answer
}
