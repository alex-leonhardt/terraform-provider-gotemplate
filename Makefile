# Check for required command tools to build or stop immediately
EXECUTABLES = git go dep find pwd terraform
K := $(foreach exec,$(EXECUTABLES),\
        $(if $(shell which $(exec)),some string,$(error "No $(exec) in PATH)))

ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

BINARY=terraform-provider-gotemplate
VERSION=1.0.0
BUILD=`git rev-parse HEAD`
PLATFORMS=darwin linux windows
ARCHITECTURES=386 amd64

# Setup linker flags option for build that interoperate with variable names in src code
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD}"

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
		UNAME_P := $(shell uname -p)
	ifeq ($(UNAME_P),x86_64)
		GOARCH = amd64
	endif
		ifneq ($(filter %86,$(UNAME_P)),)
			GOARCH = 386
		endif
	ifneq ($(filter arm%,$(UNAME_P)),)
		GOOS = arm
	endif
endif

default: test

all: clean build_all install

build:
	mkdir -p .terraform/plugins/$(GOOS)_$(GOARCH)
	go build ${LDFLAGS} -o .terraform/plugins/$(GOOS)_$(GOARCH)/$(BINARY)

test: build
	terraform init
	terraform apply

build_all:
	$(foreach GOOS, $(PLATFORMS),\
	$(foreach GOARCH, $(ARCHITECTURES), $(shell export GOOS=$(GOOS); export GOARCH=$(GOARCH); mkdir -p .terraform/plugins/$(GOOS)_$(GOARCH) && go build -v -o .terraform/plugins/$(GOOS)_$(GOARCH)/$(BINARY))))
	
install:
	go install ${LDFLAGS}

deps:
	dep ensure

# Remove only what we've created
clean:
	find ${ROOT_DIR} -name '${BINARY}[-?][a-zA-Z0-9]*[-?][a-zA-Z0-9]*' -delete

.PHONY: check clean install build_all all
