#!/bin/bash
# ******************************************************
# Author       	:	serialt 
# Email        	:	tserialt@gmail.com
# Filename     	:   makefile
# Version      	:	v1.3.0
# Created Time 	:	2021-06-25 10:47
# Last modified	:	2021-06-25 10:47
# By Modified  	: 
# Description  	:       build go package
#  
# ******************************************************


PROJECT_NAME= cli


GOBASE=$(shell pwd)
GOFILES=$(wildcard *.go)


BRANCH := $(shell git symbolic-ref HEAD 2>/dev/null | cut -d"/" -f 3)
# BRANCH := `git fetch --tags && git tag | sort -V | tail -1`
BUILD := $(shell git rev-parse --short HEAD)
BUILD_DIR := $(GOBASE)/build
VERSION = $(BRANCH)-$(BUILD)

BuildTime := $(shell date -u  '+%Y-%m-%d %H:%M:%S %Z')
GitHash := $(shell git rev-parse HEAD)
GoVersion := $(shell go version | awk '{print $3}')
Maintainer := tserialt@gmail.com 

PKGFLAGS := " -X 'main.APPName=$(PROJECT_NAME)' -X 'main.APPVersion=$(VERSION)'  -X 'main.Maintainer=$(Maintainer)'  -X 'main.BuildTime=$(BuildTime)' -X 'main.GitCommit=$(GitHash)' "

APP_NAME = $(PROJECT_NAME)-$(VERSION)
# go-pkg.v0.1.1-66ee01c-linux-amd64

.PHONY: clean
clean:
	\rm -rf build/$(PROJECT_NAME)* 

.PHONY: serve
serve:
	go run .

.PHONY: build
build: clean
	go build -ldflags $(PKGFLAGS) -o "build/$(APP_NAME)"
	@echo "编译完成"

.PHONY: release
release: clean
	go mod tidy
	GOOS="windows" GOARCH="amd64" go build -ldflags $(PKGFLAGS) -v -o "build/$(APP_NAME)-windows-amd64.exe" 
	GOOS="linux"   GOARCH="amd64" go build -ldflags $(PKGFLAGS) -v -o "build/$(APP_NAME)-linux-amd64"       
	GOOS="linux"   GOARCH="arm64" go build -ldflags $(PKGFLAGS) -v -o "build/$(APP_NAME)-linux-arm64"       
	GOOS="darwin"  GOARCH="amd64" go build -ldflags $(PKGFLAGS) -v -o "build/$(APP_NAME)-darwin-amd64"       
	GOOS="darwin"  GOARCH="arm64" go build -ldflags $(PKGFLAGS) -v -o "build/$(APP_NAME)-darwin-arm64"      
	@echo "******************"
	@echo " release succeed "
	@echo "******************"
	ls -la build/$(PROJECT_NAME)*