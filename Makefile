GO   ?= go
OS   ?= darwin 
ARCH ?= amd64

build:
	@echo ">> running build"
	@GOOS=$(OS) GOARCH=$(ARCH) $(GO) build -o proxxy . 

test:
	@echo ">> running tests"
	@$(GO) test .

format:
	@echo ">> running format"
	@$(GO) fmt .

all: format test build
