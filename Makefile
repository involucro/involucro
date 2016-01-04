
all: test lint build

PKGS = $(shell find -name \*.go | xargs dirname | uniq)

TEST_PKGS = $(shell find -name \*_test.go | xargs dirname | uniq)

get-deps:
	go get ./...
	go get github.com/smartystreets/goconvey
	go get -u github.com/golang/lint/golint

test:
	@echo Run test...
	@$(foreach pkg,$(TEST_PKGS),go test -v $(pkg) || exit;)

build:
	@CGO_ENABLED=0 go build -ldflags "-s -w" ./.

build-win32:
	@CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o involucro32.exe ./.

build-win64:
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o involucro.exe ./.

run-convey:
	$$GOPATH/bin/goconvey -host=0.0.0.0

lint:
	@$(foreach pkg,$(PKGS),$$GOPATH/bin/golint $(pkg) || exit;)

.PHONY: test build run get-deps all run-convey lint
