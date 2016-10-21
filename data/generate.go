package data
//go:generate git submodule init
//go:generate git submodule update
//go:generate go run generate_classifier.go 
//go:generate go-bindata -pkg data -o data.go classifier
