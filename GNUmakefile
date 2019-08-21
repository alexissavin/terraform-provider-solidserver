SHELL := /bin/bash
GO_FILES?=$(find . -name '*.go' |grep -v vendor)
PKG_NAME=solidserver

ifdef VERSION
	RELEASE := v$(VERSION)
else
	RELEASE := latest
endif

default: build

build:
	go get -v ./...
	go mod tidy
	go mod vendor
	if ! [ -d './_test' ]; then mkdir './_test'; fi
	go build -o ./_test/terraform-provider-solidserver

release:
	go get -v ./...
	go mod tidy
	go mod vendor
	if ! [ -d './_releases' ]; then mkdir './_release'; fi
	if ! [ -d "./_releases/$(RELEASE)" ]; then mkdir "./_releases/$(RELEASE)"; else rm -rf ./_releases/$(RELEASE)/*; fi
	if ! [ -d "./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)-linux_amd64" ]; then mkdir "./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)-linux_amd64"; fi
	if ! [ -d "./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)-freebsd_amd64" ]; then mkdir "./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)-freebsd_amd64"; fi
	if ! [ -d "./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)-windows_amd64" ]; then mkdir "./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)-windows_amd64"; fi
	cp -r ./README.md ./USAGE.md ./LICENSE ./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)-linux_amd64/
	cp -r ./README.md ./USAGE.md ./LICENSE ./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)-freebsd_amd64/
	cp -r ./README.md ./USAGE.md ./LICENSE ./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)-windows_amd64/
	env GOOS=linux GOARCH=amd64 go build -o ./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)-linux_amd64/terraform-provider-solidserver_$(RELEASE)-linux_amd64
	env GOOS=freebsd GOARCH=amd64 go build -o ./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)-freebsd_amd64/terraform-provider-solidserver_$(RELEASE)-freebsd_amd64
	env GOOS=windows GOARCH=amd64 go build -o ./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)-windows_amd64/terraform-provider-solidserver_$(RELEASE)-windows_amd64
	zip -j -r ./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)-linux_amd64.zip ./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)-linux_amd64
	zip -j -r ./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)-freebsd_amd64.zip ./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)-freebsd_amd64
	zip -j -r ./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)-windows_amd64.zip ./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)-windows_amd64
	rm -rf ./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)-linux_amd64
	rm -rf ./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)-freebsd_amd64
	rm -rf ./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)-windows_amd64

test: fmtcheck vet
	go test -v ./... || exit 1

fmt:
	gofmt -s -w ./*.go
	gofmt -s -w ./solidserver/*.go

vet:
	go vet -all ./solidserver

fmtcheck:
	./scripts/gofmtcheck.sh
