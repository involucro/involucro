
all: test

SOURCES = $(wildcard *.go)
PKGS = ./.

get-deps:
	go get ./...

test:
	$(foreach pkg,$(PKGS),go test $(pkg) || exit;)

build:
	@go build ./.

run:
	@go run $(SOURCES)
