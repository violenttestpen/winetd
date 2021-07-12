CC = gcc

CCFLAGS = -g

SVC_DIRNAME = ./servers
SVC_FILENAME = $(SVC_DIRNAME)/echo.c
SVC_DIST_FILENAME = $(SVC_DIRNAME)/echo.exe

GO = go
GOBUILD = $(GO) build
GORUN = $(GO) run
GOTEST = $(GO) test
GOBENCH = $(GOTEST) -bench=. -benchmem

BUILDFLAGS = -ldflags="-s -w"

SRC_FILES = .
DIST_DIRNAME = ./dist
DIST_FILEPATH = $(DIST_DIRNAME)/winetd.exe

all: $(DIST_FILEPATH)

mkdirs:
	mkdir $(DIST_DIRNAME) | echo

$(DIST_FILEPATH): mkdirs
	$(GOBUILD) -o $(DIST_FILEPATH) $(BUILDFLAGS) $(SRC_FILES)

services:
	$(CC) -o $(SVC_DIST_FILENAME) $(CCFLAGS) $(SVC_FILENAME)

run:
	$(GORUN) $(SRC_FILES) -server $(SVC_DIST_FILENAME) -verbosity 2 -bind 127.0.0.1

run-pwn:
	$(GORUN) $(SRC_FILES) -server $(SVC_DIRNAME)/vuln.exe -verbosity 2 -bind 127.0.0.1

benchmark:
	$(GOBENCH)

clean: $(DIST_DIRNAME)
	rm -rf $(DIST_DIRNAME)
