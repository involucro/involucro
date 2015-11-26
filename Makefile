
all: test lint

SOURCES = $(wildcard **/*.go)
PKGS = ./.

get-deps:
	@go get ./...
	@go get github.com/smartystreets/goconvey
	@go get -u github.com/golang/lint/golint

test:
	@echo Run test...
	@$(foreach pkg,$(PKGS),go test -v $(pkg) || exit;)

build:
	@go build ./.

run:
	@go run $(SOURCES)

run-convey:
	$$GOPATH/bin/goconvey

lint:
	@$$GOPATH/bin/golint

.PHONY: test build run get-deps all run-convey lint
