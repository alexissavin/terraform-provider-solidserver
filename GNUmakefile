SHELL := /bin/bash
GO_FILES?=$(find . -name '*.go' |grep -v vendor)
PKG_NAME=solidserver

default: build

build:
	go get -t -v ./...
	if ! [ -d './_test' ]; then mkdir './_test'; fi
	go build -o ./_test/terraform-provider-solidserver

test: fmtcheck vet
	export TF_ACC=1 && go test -v ./... || exit 1

fmt:
	gofmt -s -w ./*.go
	gofmt -s -w ./solidserver/*.go

vet:
	go vet -all ./solidserver

fmtcheck:
	./scripts/gofmtcheck.sh
