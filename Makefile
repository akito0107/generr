NAME := dicon
VERSION := $(shell git tag -l | tail -1)
REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS := -X 'main.version=$(VERSION)' -X 'main.revision=$(REVISION)'
PACKAGENAME := github.com/hiko1129/generr

.PHONY: setup test main clean install

all: gen build

gen:
	go generate

build:
	go build -ldflags "$(LDFLAGS)" -o bin/generr cmd/generr/main.go

## Install dependencies
setup:
	go get -u github.com/golang/dep/cmd/dep

test:
	go test -v .

test/cover:
	go test -v -coverprofile=out .

install:
	go install $(PACKAGENAME)/cmd/generr

## remove build files
clean:
	rm -rf ./bin/*

