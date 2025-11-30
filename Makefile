.PHONY: build test clean install build-all

BINARY=sdeconvert
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS=-ldflags "-X main.Version=${VERSION}"

build:
	go build ${LDFLAGS} -o bin/${BINARY} ./cmd/sdeconvert

test:
	go test -v ./...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

clean:
	rm -rf bin/ output/ coverage.out coverage.html

install:
	go install ${LDFLAGS} ./cmd/sdeconvert

# Cross-compilation
build-all:
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o bin/${BINARY}-linux-amd64 ./cmd/sdeconvert
	GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o bin/${BINARY}-darwin-amd64 ./cmd/sdeconvert
	GOOS=darwin GOARCH=arm64 go build ${LDFLAGS} -o bin/${BINARY}-darwin-arm64 ./cmd/sdeconvert
	GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o bin/${BINARY}-windows-amd64.exe ./cmd/sdeconvert

# Development helpers
fmt:
	go fmt ./...

vet:
	go vet ./...

lint: fmt vet

tidy:
	go mod tidy
