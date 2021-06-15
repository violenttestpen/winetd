CC = gcc

CCFLAGS = -g

SVC_DIRNAME = ./servers
SVC_FILENAME = ./echo.c
SVC_DIST_FILENAME = ./echo.exe

GO = go
GOBUILD = $(GO) build
GORUN = $(GO) run
GOTEST = $(GO) test
GOBENCH = $(GOTEST) -bench=.

BUILDFLAGS = -ldflags="-s -w"

SRC_FILES = .
DIST_DIRNAME = ./dist
DIST_FILENAME = ./ginetd.exe

all: build

mkdirs:
	mkdir $(DIST_DIRNAME) | echo

build: mkdirs
	$(GOBUILD) -o $(DIST_DIRNAME)/$(DIST_FILENAME) $(BUILDFLAGS) $(SRC_FILES)

services:
	$(CC) -o $(SVC_DIRNAME)/$(SVC_DIST_FILENAME) $(CCFLAGS) $(SVC_DIRNAME)/$(SVC_FILENAME)

run:
	$(GORUN) $(SRC_FILES) -server $(SVC_DIRNAME)/$(SVC_DIST_FILENAME) -verbosity 2 -bind 127.0.0.1

run-pwn:
	$(GORUN) $(SRC_FILES) -server $(SVC_DIRNAME)/vuln.exe -verbosity 2 -bind 127.0.0.1

benchmark:
	$(GOBENCH)

clean: $(DIST_DIRNAME)
	rm -rf $(DIST_DIRNAME)
