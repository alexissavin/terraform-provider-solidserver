SHELL := /bin/bash
GO_FILES?=$(find . -name '*.go' |grep -v vendor)
PKG_NAME=solidserver

default: build

build:
	if ! [ -d './_test' ]; then mkdir './_test'; fi
	go build -o ./_test/terraform-provider-solidserver

test: fmtcheck vet
	go test -v ./... || exit 1

fmt:
	gofmt -w $(GO_FILES)

vet:
	go tool vet &> /dev/null ; if [ $$? -eq 3 ]; then go get golang.org/x/tools/cmd/vet; fi
	go tool vet -all ./solidserver

fmtcheck:
	./scripts/gofmtcheck.sh