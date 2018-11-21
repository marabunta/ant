.PHONY: all test build cover

GO ?= go
VERSION=$(shell git describe --tags --always)

build:
	${GO} build -ldflags "-s -w -X main.version=${VERSION}" -o ant cmd/ant/main.go;
	# env GOOS=freebsd GOARCH=amd64 ${GO} build -ldflags "-s -w -X main.version=${VERSION}"

test:
	${GO} test -v

clean:
	@rm -rf ant *.out

cover:
	${GO} test -cover && \
	${GO} test -coverprofile=coverage.out  && \
	${GO} tool cover -html=coverage.out
