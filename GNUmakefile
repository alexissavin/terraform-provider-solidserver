SHELL := /bin/bash
GO_FILES?=$(find . -name '*.go' |grep -v vendor)

# To provide the version use 'make release VERSION=1.1.1'
ifdef VERSION
	RELEASE := v$(VERSION)
else
	VERSION := 99999.9
	RELEASE := v99999.9
endif

# Terraform 13 local registry handler
PKG_NAME := solidserver
OS_ARCH := linux_amd64
TERRAFORM_PLUGINS_DIRECTORY := ~/.terraform.d/plugins/terraform.efficientip.com/efficientip/${PKG_NAME}/${VERSION}/${OS_ARCH}

default: build

build:
	go get -v ./...
	go mod tidy
	go mod vendor
	if ! [ -d ${TERRAFORM_PLUGINS_DIRECTORY} ]; then mkdir -p ${TERRAFORM_PLUGINS_DIRECTORY}; fi
	go build -o ${TERRAFORM_PLUGINS_DIRECTORY}/terraform-provider-${PKG_NAME}
	if [ -d ./_tests ]; then cd _tests && rm -rf .terraform* && cd ..; fi

release:
	go get -v ./...
	go mod tidy
	go mod vendor
	if ! [ -d './_releases' ]; then mkdir './_release'; fi
	if ! [ -d "./_releases/$(RELEASE)" ]; then mkdir "./_releases/$(RELEASE)"; else rm -rf ./_releases/$(RELEASE)/*; fi
	if ! [ -d "./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)-linux_amd64" ]; then mkdir "./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)-linux_amd64"; fi
	if ! [ -d "./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)-freebsd_amd64" ]; then mkdir "./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)-freebsd_amd64"; fi
	if ! [ -d "./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)-windows_amd64" ]; then mkdir "./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)-windows_amd64"; fi
	if ! [ -d "./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)-macos_amd64" ]; then mkdir "./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)-macos_amd64"; fi
	cp -r ./README.md ./USAGE.md ./LICENSE ./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)-linux_amd64/
	cp -r ./README.md ./USAGE.md ./LICENSE ./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)-freebsd_amd64/
	cp -r ./README.md ./USAGE.md ./LICENSE ./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)-windows_amd64/
	cp -r ./README.md ./USAGE.md ./LICENSE ./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)-macos_amd64/
	env GOOS=linux GOARCH=amd64 go build -o ./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)-linux_amd64/terraform-provider-solidserver_$(RELEASE)
	env GOOS=freebsd GOARCH=amd64 go build -o ./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)-freebsd_amd64/terraform-provider-solidserver_$(RELEASE)
	env GOOS=windows GOARCH=amd64 go build -o ./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)-windows_amd64/terraform-provider-solidserver_$(RELEASE)
	env GOOS=darwin GOARCH=amd64 go build -o ./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)-macos_amd64/terraform-provider-solidserver_$(RELEASE)
	zip -j -r ./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)-linux_amd64.zip ./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)-linux_amd64
	zip -j -r ./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)-freebsd_amd64.zip ./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)-freebsd_amd64
	zip -j -r ./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)-windows_amd64.zip ./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)-windows_amd64
	zip -j -r ./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)-macos_amd64.zip ./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)-macos_amd64
	rm -rf ./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)-linux_amd64
	rm -rf ./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)-freebsd_amd64
	rm -rf ./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)-windows_amd64
	rm -rf ./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)-macos_amd64

test: fmtcheck vet
	go test -v ./... || exit 1

fmt:
	gofmt -s -w ./*.go
	gofmt -s -w ./solidserver/*.go

vet:
	go vet -all ./solidserver

fmtcheck:
	./scripts/gofmtcheck.sh
