all:
	go build -o repository ./cmd/repository

.PHONY: example
example:
	cd example && go generate .

.PHONY: test
test:
	go test -v ./...
