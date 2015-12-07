package linguist

import (
	"bytes"
	"encoding/base64"
	"regexp"

	"github.com/jbrukh/bayesian"
)

var classifier *bayesian.Classifier
var classifier_initialized bool = false

func GetClassifier() *bayesian.Classifier {
	if !classifier_initialized {
		serialized, err := base64.StdEncoding.DecodeString(b64_encoded_serialized_classifier)
		if err != nil {
			panic(err.Error())
		}

		r := bytes.NewReader(serialized)
		classifier, err = bayesian.NewClassifierFromReader(r)
		if err != nil {
			panic(err.Error())
		}
		classifier_initialized = true
	}
	return classifier
}

func Analyse(b []byte) (language string) {
	// TODO(tso): port github/linguist tokenizer proper
	// instead of this copout:
	re := regexp.MustCompile(`\s+`)
	document := re.Split(string(b), -1)
	classifier := GetClassifier()
	_, id, _ := classifier.LogScores(document)
	if language, ok := class_map[id]; ok {
		return language
	}
	return ""
}
