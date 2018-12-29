# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=go-chat-app
BINARY_UNIX=$(BINARY_NAME)_unix

build:
	$(GOBUILD) -o $(BINARY_NAME) -v

run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)

test:
	$(GOTEST) -v ./...
