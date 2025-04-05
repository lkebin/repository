all:
	go build -o repository ./cmd/repository

.PHONY: testdata
testdata:
	cd testdata && go generate .
