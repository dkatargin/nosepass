.GOOS ?=
GOARCH ?=
GO111MODULE ?= on
GOPATH ?= $(CURDIR)

ifeq ($(GOOS),)
  GOOS = $(shell go version | awk -F ' ' '{print $$NF}' | awk -F '/' '{print $$1}')
endif
ifeq ($(GOARCH),)
  GOARCH = $(shell go version | awk -F ' ' '{print $$NF}' | awk -F '/' '{print $$2}')
endif
ifeq ($(VERSION),)
  VERSION = latest
endif

PACKAGES = $(shell $(GO) list ./... | grep -v '/vendor/')
GO := GOPATH=$(GOPATH) GO111MODULE=$(GO111MODULE) go

.DEFAULT_GOAL := build

.PHONY: show_env
show_env:
	@echo ">> show env"
	@echo "   GOOS              = $(GOOS)"
	@echo "   GOARCH            = $(GOARCH)"
	@echo "   GOROOT            = $(GOROOT)"
	@echo "   GOPATH            = $(GOPATH)"
	@echo "   GO111MODULE       = $(GO111MODULE)"
	@echo "   VERSION           = $(VERSION)"
	@echo "   PACKAGES          = $(PACKAGES)"


.PHONY: build
build: show_env
	@echo ">> build sociald"
	$(GO) build -o ./bin/mypass .
	@echo ">> done"

.PHONY: clean
clean:
	@echo ">> cleanup"
	rm -rf ./bin/
	rm -rf ./src/pkg/
	@echo ">> done"
