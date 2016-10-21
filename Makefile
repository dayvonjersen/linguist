all:
	go generate
	go generate ./data
	go install
	go install ./cmd/l
