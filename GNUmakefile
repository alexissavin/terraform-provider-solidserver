TEST?=$$(go list ./... |grep -v 'vendor')
TESTARGS?=''
GO_FILES?=$$(find . -name '*.go' |grep -v vendor)
WEBSITE_REPO=github.com/hashicorp/terraform-website
PKG_NAME=template

default: build

build:
	if ! [ -d './_test' ]; then mkdir './_test'; fi
	@go build -o ./_test/terraform-provider-solidserver

test: fmtcheck
	@go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

fmt:
	@gofmt -w $(GO_FILES)

vet:
	@go tool vet 2>/dev/null ; if [ $$? -eq 3 ]; then \
		go get golang.org/x/tools/cmd/vet; \
	fi
	echo "Running 'go tool vet $(VETARGS) $(GO_FILES)'"
	@go tool vet $(VETARGS) $(GO_FILES) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Suspicious constructs found, please fix it."; \
		exit 1; \
	fi

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"
