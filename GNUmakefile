SHELL := /bin/bash
GO_FILES?=$(find . -name '*.go' |grep -v vendor)

# To provide the version use 'make release VERSION=1.1.1 GPGKEY=<example@efficientip.com>'
ifdef VERSION
	RELEASE := $(VERSION)
else
	RELEASE := 99999.9
endif

ifdef GPGKEY
	GPGKEYOPTION := -u $(GPGKEY)
else
	GPGKEYOPTION :=
endif

# Terraform 13 local registry handler
PKG_NAME := solidserver
OS_ARCH := linux_amd64
TERRAFORM_PLUGINS_DIRECTORY := ~/.terraform.d/plugins/terraform.efficientip.com/efficientip/${PKG_NAME}/${RELEASE}/${OS_ARCH}

default: build

build:
	go get -v ./...
	go mod tidy
	go mod vendor
	if ! [ -d ${TERRAFORM_PLUGINS_DIRECTORY} ]; then mkdir -p ${TERRAFORM_PLUGINS_DIRECTORY}; fi
	go build -o ${TERRAFORM_PLUGINS_DIRECTORY}/terraform-provider-${PKG_NAME}
	if [ -d ./_tests ]; then cd _tests && rm -rf .terraform* && cd ..; fi

release:
  #The binary name format is terraform-provider-{NAME}_v{VERSION}
  #The archive name format is terraform-provider-{NAME}_{VERSION}_{OS}_{ARCH}.zip
	go get -v ./...
	go mod tidy
	go mod vendor
	if ! [ -d './_releases' ]; then mkdir './_release'; fi
	if ! [ -d "./_releases/$(RELEASE)" ]; then mkdir "./_releases/$(RELEASE)"; else rm -rf ./_releases/$(RELEASE)/*; fi
	if ! [ -d "./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)_linux_amd64" ]; then mkdir "./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)_linux_amd64"; fi
	if ! [ -d "./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)_freebsd_amd64" ]; then mkdir "./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)_freebsd_amd64"; fi
	if ! [ -d "./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)_windows_amd64" ]; then mkdir "./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)_windows_amd64"; fi
	if ! [ -d "./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)_darwin_amd64" ]; then mkdir "./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)_darwin_amd64"; fi
	if ! [ -d "./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)_darwin_arm64" ]; then mkdir "./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)_darwin_arm64"; fi
	cp -r ./README.md ./LICENSE ./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)_linux_amd64/
	cp -r ./README.md ./LICENSE ./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)_freebsd_amd64/
	cp -r ./README.md ./LICENSE ./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)_windows_amd64/
	cp -r ./README.md ./LICENSE ./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)_darwin_amd64/
	cp -r ./README.md ./LICENSE ./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)_darwin_arm64/
	env GOOS=linux GOARCH=amd64 go build -o ./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)_linux_amd64/terraform-provider-solidserver_v$(RELEASE)
	env GOOS=freebsd GOARCH=amd64 go build -o ./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)_freebsd_amd64/terraform-provider-solidserver_v$(RELEASE)
	env GOOS=windows GOARCH=amd64 go build -o ./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)_windows_amd64/terraform-provider-solidserver_v$(RELEASE)
	env GOOS=darwin GOARCH=amd64 go build -o ./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)_darwin_amd64/terraform-provider-solidserver_v$(RELEASE)
	env GOOS=darwin GOARCH=arm64 go build -o ./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)_darwin_arm64/terraform-provider-solidserver_v$(RELEASE)
	zip -j -r ./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)_linux_amd64.zip ./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)_linux_amd64
	zip -j -r ./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)_freebsd_amd64.zip ./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)_freebsd_amd64
	zip -j -r ./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)_windows_amd64.zip ./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)_windows_amd64
	zip -j -r ./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)_darwin_amd64.zip ./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)_darwin_amd64
	zip -j -r ./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)_darwin_arm64.zip ./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)_darwin_arm64
	rm -rf ./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)_linux_amd64
	rm -rf ./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)_freebsd_amd64
	rm -rf ./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)_windows_amd64
	rm -rf ./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)_darwin_amd64
	rm -rf ./_releases/$(RELEASE)/terraform-provider-solidserver_$(RELEASE)_darwin_arm64
	cd ./_releases/$(RELEASE) && shasum -a 256 *.zip > ./terraform-provider-solidserver_$(RELEASE)_SHA256SUMS && cd ../..
	cd ./_releases/$(RELEASE) && gpg $(GPGKEYOPTION) --detach-sign ./terraform-provider-solidserver_$(RELEASE)_SHA256SUMS && cd ../..

test: fmtcheck vet
	go test -v ./... || exit 1

fmt:
	gofmt -s -w ./*.go
	gofmt -s -w ./solidserver/*.go

vet:
	go vet -all ./solidserver

fmtcheck:
	./scripts/gofmtcheck.sh