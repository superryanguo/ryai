SHELL:=/bin/bash
.PHONY: ryai test clean

export GO111MODULE=on

# Path Related
MKFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
MKFILE_DIR := $(dir $(MKFILE_PATH))
RELEASE_DIR := ${MKFILE_DIR}bin
GO_PATH := $(shell go env | grep GOPATH | awk -F '"' '{print $$2}')
VERSION=$(shell git describe --tags --dirty || echo "unknown version")
BDTIME=$(shell date --utc "+%Y-%m-%d %H:%M:%S" || echo "unknown date")
GOBUILD=go build -v -trimpath -ldflags '-X "github.com/superryanguo/ryai/utils.Version=$(VERSION)" \
		-X "github.com/superryanguo/ryai/utils.BuildTime=$(BDTIME)" \
		-w -s -buildid=$(VERSION)'

# Image Name
IMAGE_NAME?=ryai

# Version
RELEASE?=v0.1

# IP
IP?=192.168.1.1
# Git Related
GIT_REPO_INFO=$(shell cd ${MKFILE_DIR} && git config --get remote.origin.url)
ifndef GIT_COMMIT
  GIT_COMMIT := git-$(shell git rev-parse --short HEAD)
endif

# Rules
ryai: main.go
	cd ${MKFILE_DIR} && \
	GOOS=linux $(GOBUILD) -o $(RELEASE_DIR)/$@

test:
	cd ${MKFILE_DIR}
	go mod tidy
	git diff --exit-code go.mod go.sum
	go mod verify
	go test -v -gcflags "all=-l" ${MKFILE_DIR}pkg/... ${TEST_FLAGS}

clean:
	rm -rf ${RELEASE_DIR}

run: ryai
	${RELEASE_DIR}/ryai
