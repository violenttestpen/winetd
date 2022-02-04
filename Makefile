CC = gcc

CCFLAGS = -g

SVC_DIRNAME = ./servers
SVC_FILENAME = $(SVC_DIRNAME)/echo.c
SVC_DIST_FILENAME = $(SVC_DIRNAME)/echo.exe

GO = go
GOBUILD = $(GO) build
GORUN = $(GO) run
GOTEST = $(GO) test
GOBENCH = $(GOTEST) -run=none -bench=. -benchmem

BUILDFLAGS = -ldflags="-s -w"

SRC_FILES = .
DIST_FILEPATH = ./winetd.exe

all: $(DIST_FILEPATH)

$(DIST_FILEPATH):
	$(GOBUILD) -o $(DIST_FILEPATH) $(BUILDFLAGS) $(SRC_FILES)

services:
	$(CC) -o $(SVC_DIST_FILENAME) $(CCFLAGS) $(SVC_FILENAME)

run:
	$(GORUN) -race $(SRC_FILES) -server $(SVC_DIST_FILENAME) -verbosity 2 -bind 127.0.0.1

run-pwn:
	$(GORUN) -race $(SRC_FILES) -server $(SVC_DIRNAME)/vuln.exe -verbosity 2 -bind 127.0.0.1

benchmark:
	$(GOBENCH) -v

clean:
	rm -rf $(DIST_FILEPATH) 2>/dev/null
