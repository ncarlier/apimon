.SILENT :

# Author
AUTHOR=github.com/ncarlier

# App name
APPNAME=apimon

# Go configuration
GOOS?=linux
GOARCH?=amd64

# Add exe extension if windows target
is_windows:=$(filter windows,$(GOOS))
EXT:=$(if $(is_windows),".exe","")

# Go app path
APPBASE=${GOPATH}/src/$(AUTHOR)

# Artefact name
ARTEFACT=release/$(APPNAME)-$(GOOS)-$(GOARCH)$(EXT)

# Extract version infos
VERSION:=`git describe --tags`
LDFLAGS=-ldflags "-X $(AUTHOR)/$(APPNAME)/main.Version=${VERSION}"

all: build

# Include common Make tasks
root_dir:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
makefiles:=$(root_dir)/makefiles
include $(makefiles)/help.Makefile

$(APPBASE)/$(APPNAME):
	echo "Creating GO src link: $(APPBASE)/$(APPNAME) ..."
	mkdir -p $(APPBASE)
	ln -s $(root_dir) $(APPBASE)/$(APPNAME)

## Clean built files
clean:
	-rm -rf release
.PHONY: clean

deps:
	echo ">>> Installing dependencies ..."
	cd $(APPBASE)/$(APPNAME) && dep ensure
.PHONY: deps

Gopkg.lock:
	make deps

## Build executable
build: Gopkg.lock $(APPBASE)/$(APPNAME)
	-mkdir -p release
	echo "Building: $(ARTEFACT) ..."
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(LDFLAGS) -o $(ARTEFACT)
.PHONY: build

$(ARTEFACT): build

## Run tests
test:
	cd $(APPBASE)/$(APPNAME) && go test `go list ./... | grep -v vendor`
.PHONY: test

## Install executable
install: $(ARTEFACT)
	echo "Installing $(ARTEFACT) to ${HOME}/.local/bin/$(APPNAME) ..."
	cp $(ARTEFACT) ${HOME}/.local/bin/$(APPNAME)
.PHONY: install

## Create Docker image
image:
	echo "Building Docker inage ..."
	docker build --rm -t ncarlier/$(APPNAME) .
.PHONY: image

## Generate changelog
changelog:
	standard-changelog --first-release
.PHONY: changelog

## GZIP executable
gzip:
	tar cvzf $(ARTEFACT).tgz $(ARTEFACT)
.PHONY: gzip

## Create distribution binaries
distribution:
	GOARCH=amd64 make build gzip
	GOARCH=arm64 make build gzip
	GOARCH=arm make build gzip
	GOOS=darwin make build gzip
	GOOS=windows make build gzip
.PHONY: distribution
