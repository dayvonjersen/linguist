package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jbrukh/bayesian"
)

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	log.SetOutput(os.Stderr)
	log.Println("Opening ../samples.json ...")
	f, err := os.Open("../samples.json")
	checkErr(err)
	log.Println("Reading ../samples.json ...")
	b, err := ioutil.ReadAll(f)
	checkErr(err)

	log.Println("Unmarshaling ../samples.json ...")
	var dat map[string]interface{}
	checkErr(json.Unmarshal(b, &dat))

	tokens := dat["tokens"].(map[string]interface{})

	classes := make([]bayesian.Class, 1)
	documents := make(map[bayesian.Class][]string)
	results := make(map[int]string)
	index := 0

	log.Println("Reading data into vars ...")
	for key, token := range tokens {
		if key == "" {
			continue
		}

		results[index] = key
		index++

		key = strings.Replace(key, " ", "_", -1)
		key = strings.Replace(key, "#", "SHARP", -1)
		key = strings.Replace(key, "+", "PLUS", -1)
		key = strings.Replace(key, "-", "DASH", -1)
		var class bayesian.Class = bayesian.Class(key)

		classes = append(classes, class)

		words := token.(map[string]interface{})
		w := make([]string, 1)
		for word, repeat := range words {
			count := repeat.(float64)
			for r := 0.0; r < count; r++ {
				w = append(w, word)
			}
		}
		documents[class] = w
	}

	log.Println("Generating go code ...")

	fmt.Println("package linguist")
	fmt.Println("\n// *** DO NOT MODIFY THIS FILE DIRECTLY *** \\\\")
	fmt.Println("// this file was automatically generated ")
	fmt.Println("// by generate_classifier.go at")
	fmt.Printf("// %s\n\n", time.Now())
    fmt.Printf("var class_map map[int]string = %#v\n", results)

	log.Println("Creating bayesian.Classifier ...")
	clsf := bayesian.NewClassifier(classes...)
	for cls, dox := range documents {
		clsf.Learn(dox, cls)
	}

	log.Println("Serializing and exporting bayesian.Classifier into ./out/classifier ...")
	checkErr(clsf.WriteToFile("./out/classifier"))

	var buf bytes.Buffer
	b64 := base64.NewEncoder(base64.StdEncoding, &buf)
	checkErr(clsf.WriteTo(b64))

	log.Println("Generating more go code...")
	b64_data, err := ioutil.ReadAll(&buf)
	checkErr(err)
	fmt.Printf("var b64_encoded_serialized_classifier string = `%s`", string(b64_data))

	log.Println("Done.")
}
