#!/usr/bin/env bash
set -euo pipefail

PKGS=(./internal/... ./auth)
VERSION="$(git describe --tags --always)"
LDFLAGS="-s -w -X github.com/involucro/involucro/app.version=$VERSION"

case "$MODE" in
unit)
    go test -v -short "${PKGS[@]}"
    CGO_ENABLED=0 go build -o "$FILENAME" -ldflags "$LDFLAGS" ./cmd/involucro
    "./$FILENAME" --version
    ;;

integration)
    CGO_ENABLED=0 go build -o "$FILENAME" -ldflags "$LDFLAGS" ./cmd/involucro
    "./$FILENAME" --version
    "./$FILENAME" wrap-yourself
    go test -v "${PKGS[@]}"
    ;;

windows-amd64)
    mkdir -p .bin
    env GOFLAGS= GOBIN="$PWD/.bin" go install github.com/josephspurrier/goversioninfo/cmd/goversioninfo@v1.7.0
    ./.bin/goversioninfo "-product-version=$VERSION"

    CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o involucro.exe -ldflags "$LDFLAGS" ./cmd/involucro
    file involucro.exe
    ;;

linux-amd64)
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o "$FILENAME" -ldflags "$LDFLAGS" ./cmd/involucro
    file "$FILENAME"
    ;;

linux-arm64)
    CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o "$FILENAME" -ldflags "$LDFLAGS" ./cmd/involucro
    file "$FILENAME"
    ;;

linux-arm-v7)
    CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -o "$FILENAME" -ldflags "$LDFLAGS" ./cmd/involucro
    file "$FILENAME"
    ;;

linux-riscv64)
    CGO_ENABLED=0 GOOS=linux GOARCH=riscv64 go build -o "$FILENAME" -ldflags "$LDFLAGS" ./cmd/involucro
    file "$FILENAME"
    ;;

*)
    echo "Unknown MODE: $MODE" >&2
    exit 1
    ;;
esac
