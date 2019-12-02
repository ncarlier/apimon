.SILENT :

export GO111MODULE=on

# App name
APPNAME=apimon

# Go configuration
GOOS?=linux
GOARCH?=amd64

# Add exe extension if windows target
is_windows:=$(filter windows,$(GOOS))
EXT:=$(if $(is_windows),".exe","")

# Archive name
ARCHIVE=$(APPNAME)-$(GOOS)-$(GOARCH).tgz

# Executable name
EXECUTABLE=$(APPNAME)$(EXT)

# Extract version infos
# Extract version infos
VERSION:=`git describe --tags`
GIT_COMMIT:=`git rev-list -1 HEAD`
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.GitCommit=${GIT_COMMIT}"

all: build

# Include common Make tasks
root_dir:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
makefiles:=$(root_dir)/makefiles
include $(makefiles)/help.Makefile
include $(makefiles)/docker/compose.Makefile

## Clean built files
clean:
	echo ">>> Cleaning..."
	-rm -rf release
.PHONY: clean

## Build executable
build:
	-mkdir -p release
	echo ">>> Building: $(EXECUTABLE) ..."
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(LDFLAGS) -o release/$(EXECUTABLE)
.PHONY: build

release/$(EXECUTABLE): build

## Run tests
test:
	echo ">>> Running tests..."
	go test ./...
.PHONY: test

## Install executable
install: release/$(EXECUTABLE)
	echo ">>> Installing $(EXECUTABLE) to ${HOME}/.local/bin/$(EXECUTABLE) ..."
	install release/$(EXECUTABLE) ~/.local/bin/
.PHONY: install

## Create Docker image
image:
	echo ">>> Building Docker image..."
	docker build --rm -t ncarlier/$(APPNAME) .
.PHONY: image

## Generate changelog
changelog:
	standard-changelog --first-release
.PHONY: changelog

## Create archive
archive:
	echo ">>> Creating release/$(ARCHIVE) archive..."
	tar czf release/$(ARCHIVE) \
		--exclude=*.tgz \
	 	README.md \
		LICENSE \
		CHANGELOG.md \
		-C release/ $(subst release/,,$(wildcard release/*))
	find release/ -type f -not -name '*.tgz' -delete
.PHONY: archive

## Create distribution binaries
distribution:
	GOARCH=amd64 make build archive
	GOARCH=arm64 make build archive
	GOARCH=arm   make build archive
	GOOS=windows make build archive
	GOOS=darwin  make build archive
.PHONY: distribution

## Deploy Docker stack
deploy: compose-up 
.PHONY: deploy

## Un-deploy Docker stack
undeploy: compose-down
.PHONY: undeploy
