

GO           ?= go
GOFMT        ?= $(GO)fmt
FIRST_GOPATH := $(firstword $(subst :, ,$(shell $(GO) env GOPATH)))

BIN_DIR ?= $(shell pwd)/build
VERSION ?= $(shell cat VERSION)


all:  fmt style  cpimigrator 
 
style:
	@echo ">> checking code style"
	! $(GOFMT) -d $$(find . -name '*.go' -print) | grep '^'

cpimigrator: 
	@echo ">> building cpimigrator binaries"
	$(GO) build -o build/cpimigrator cmd/cpimigrator/main.go

fmt:
	@echo ">> format code style"
	$(GOFMT) -w $$(find . -name '*.go' -print) 

clean:
	rm -rf $(BIN_DIR)

.PHONY: all style  fmt  cpimigrator