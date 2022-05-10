GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GORUN=$(GOCMD) run
GOCLEAN=$(GOCMD) clean
GOINSTALL=$(GOCMD) install
BINARY_NAME=unityguid
BINARY_DIR=bin

all: clean build system_install

.PHONY: deps
deps:
	$(GOCMD) get

.PHONY: deps-test
deps-test:
	$(GOCMD) get github.com/stretchr/testify/assert

.PHONY: test
test:
	./script/test.sh

.PHONY: build
build:
	$(GOBUILD) -o $(BINARY_DIR)/$(BINARY_NAME) *.go

.PHONY: run
run:
	$(GORUN) *.go

.PHONY: clean
clean:
	$(GOCLEAN)
	rm -rf $(BINARY_DIR)

.PHONY: system_install
system_install:
	$(GOINSTALL)

