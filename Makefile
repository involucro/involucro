
build:
	@CGO_ENABLED=0 go build -ldflags "-s -w" ./.

build-win32:
	@CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o involucro32.exe ./.

build-win64:
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o involucro.exe ./.
