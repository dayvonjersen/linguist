all:
	go generate
	go generate ./data
	go install
	go get ./cmd/l
