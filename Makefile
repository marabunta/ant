.PHONY: all test build clean cover

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

certs: SHELL:=/bin/bash
certs:
	@mkdir -p certs
	# crate CA
	openssl req -x509 -nodes -days 3650 -newkey ec:<(openssl ecparam -name prime256v1) -keyout certs/CA.key -out certs/CA.crt -subj "/CN=marabunta"

	# create server certs
	openssl req -new -nodes -newkey ec:<(openssl ecparam -name prime256v1) -keyout certs/server.key -out certs/server.csr -subj "/CN=marabunta.host"
	openssl x509 -days 3065 -sha256 -req -in certs/server.csr -CA certs/CA.crt -CAkey certs/CA.key -set_serial 01 -out certs/server.crt

	# create client certs
	openssl req -new -nodes -newkey ec:<(openssl ecparam -name prime256v1) -keyout certs/client.key -out certs/client.csr -subj "/CN=client.marabunta.host"
	openssl x509 -days 3065 -sha256 -req -in certs/client.csr -CA certs/CA.crt -CAkey certs/CA.key -set_serial 01 -out certs/client.crt
