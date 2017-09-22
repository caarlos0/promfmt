TEST_PATTERN?=.
TEST_OPTIONS?=

# Install all the build and lint dependencies
setup:
	go get -u github.com/alecthomas/gometalinter
	go get -u github.com/golang/dep/...
	go get -u github.com/pierrre/gotestcover
	go get -u golang.org/x/tools/cmd/cover
	dep ensure
	gometalinter --install --update

# Run all the tests
test:
	gotestcover $(TEST_OPTIONS) -covermode=count -coverprofile=coverage.out ./... -run $(TEST_PATTERN) -timeout=30s

# Run all the tests and opens the coverage report
cover: test
	go tool cover -html=coverage.out

# gofmt and goimports all go files
fmt:
	find . -name '*.go' -not -wholename './vendor/*' | while read -r file; do gofmt -w -s "$$file"; goimports -w "$$file"; done

# Run all the linters
lint:
	gometalinter --vendor ./...

# Run all the tests and code checks
ci: lint test

# build a local version
build:
	go build .

.DEFAULT_GOAL := build
