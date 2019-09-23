.PHONY: darwin linux all

darwin: Darwin
Darwin:
	cd cmd/yamlmerge && GOOS=darwin GOARCH=amd64 go build -v -o  ../../bin/yamlmerge .

linux:Linux
Linux:
	cd cmd/yamlmerge && GOOS=linux GOARCH=amd64 go build -a -ldflags '-extldflags "-static"' -v -o ../../bin/yamlmerge .

all:
	make `uname -s`