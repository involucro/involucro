
all: test

SOURCES = $(wildcard *.go)
PKGS = ./.

get-deps:
	@go get ./...
	@go get github.com/smartystreets/goconvey

test:
	@echo Run test...
	@$(foreach pkg,$(PKGS),go test -v $(pkg) || exit;)

build:
	@go build ./.

run:
	@go run $(SOURCES)

.PHONY: test build run get-deps all
