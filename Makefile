.PHONY: darwin linux all

darwin:
	cd cmd/yamlmerge && GOOS=darwin GOARCH=amd64 go build -v -o  ../../bin/darwin/yamlmerge .

linux:
	cd cmd/yamlmerge && GOOS=linux GOARCH=amd64 go build -a -ldflags '-extldflags "-static"' -v -o ../../bin/linux/yamlmerge .

all: darwin linux

