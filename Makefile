SOURCE_FILES?=./...
TEST_PATTERN?=.
TEST_OPTIONS?=

export PATH := $(PWD)/bin:$(PATH)
export GO111MODULE := on
export GOPROXY = https://proxy.golang.org,direct
export GOBIN = $(PWD)/bin

GOLANGCI_LINT := $(GOBIN)/golangci-lint
$(GOLANGCI_LINT): 
	cd "$$(mktemp -d)" && \
		go get github.com/golangci/golangci-lint/cmd/golangci-lint

# Install all the build and lint dependencies
setup: $(GOLANGCI_LINT)
.PHONY: setup

# Run all the tests
test:
	go test $(TEST_OPTIONS) -failfast -race -coverpkg=./... -covermode=atomic -coverprofile=coverage.txt $(SOURCE_FILES) -run $(TEST_PATTERN) -timeout=2m
.PHONY: test

# Run all the tests and opens the coverage report
cover: test
	go tool cover -html=coverage.txt
.PHONY: cover

# gofmt and goimports all go files
fmt:
	find . -name '*.go' -not -wholename './vendor/*' | while read -r file; do gofmt -w -s "$$file"; goimports -w "$$file"; done
.PHONY: fmt

# Run all the linters
lint: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run $(LINT_OPTIONS)
.PHONY: lint

# Run all the tests and code checks
ci: build test lint
.PHONY: ci

# build a local version
build:
	go build
.PHONY: build

.DEFAULT_GOAL := build
