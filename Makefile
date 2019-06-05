TEST?="./gotemplate"
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
PKG_NAME=gotemplate

# Detect OS
GOOS :=
GOARCH :=
ifeq ($(OS),Windows_NT)
	GOOS = windows
	ifeq ($(PROCESSOR_ARCHITECTURE),AMD64)
		GOARCH = amd64
	endif
	ifeq ($(PROCESSOR_ARCHITECTURE),x86)
		GOARCH = 386
	endif
else
	UNAME_S := $(shell uname -s)
	ifeq ($(UNAME_S),Linux)
		GOOS = linux
	endif
	ifeq ($(UNAME_S),Darwin)
		GOOS = darwin
	endif
		UNAME_M := $(shell uname -m)
	ifeq ($(UNAME_M),x86_64)
		GOARCH = amd64
	endif
		ifneq ($(filter %86,$(UNAME_M)),)
			GOARCH = 386
		endif
	ifneq ($(filter arm%,$(UNAME_M)),)
		GOOS = arm
	endif
endif

default: install

install: errcheck fmtcheck
	go install

build: errcheck fmtcheck
	$(shell '$(CURDIR)/scripts/build.sh' $(GOOS) $(GOARCH))

build_all:
	$(foreach GOOS, darwin linux windows, $(foreach GOARCH, 386 amd64, $(shell '$(CURDIR)/scripts/build.sh' $(GOOS) $(GOARCH))))

install_plugin: build
	$(shell '$(CURDIR)/scripts/install.sh' $(GOOS) $(GOARCH))

test: fmtcheck
	go test $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	$(shell '$(CURDIR)/scripts/gofmtcheck.sh')

errcheck:
	$(shell '$(CURDIR)/scripts/errcheck.sh')

lint:
	@echo "==> Checking source code against linters..."
	@GOGC=30 golangci-lint --verbose run ./$(PKG_NAME)

tools:
	GO111MODULE=on go install github.com/client9/misspell/cmd/misspell
	GO111MODULE=on go install github.com/golangci/golangci-lint/cmd/golangci-lint

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

.PHONY: build build_all test testacc vet fmt fmtcheck errcheck lint tools test-compile
