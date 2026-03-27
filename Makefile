.PHONY: build test lint clean install

VERSION ?= 1.0.0
LDFLAGS = -ldflags "-X github.com/JSLEEKR/schemadiff/cmd.Version=$(VERSION)"

build:
	go build $(LDFLAGS) -o schemadiff .

test:
	go test -v -race -count=1 ./...

lint:
	go vet ./...

cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

clean:
	rm -f schemadiff schemadiff.exe coverage.out coverage.html
	rm -rf dist/

install:
	go install $(LDFLAGS) .

dist: clean
	mkdir -p dist
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/schemadiff-linux-amd64 .
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/schemadiff-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/schemadiff-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/schemadiff-windows-amd64.exe .
