.PHONY: all

all:
	cd cmd/configgen && GOOS=darwin GOARCH=amd64 go build -v -o  ../../bin/configgen .

linux:
	cd cmd/configgen && GOOS=linux GOARCH=amd64 go build -a -ldflags '-extldflags "-static"' -v -o ../../bin/configgen .

