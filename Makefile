.PHONY: all

all:
	cd cmd/configgen && GOOS=darwin GOARCH=amd64 go build -v -o  ../../bin/yamlmerge .

linux:
	cd cmd/yamlmerge && GOOS=linux GOARCH=amd64 go build -a -ldflags '-extldflags "-static"' -v -o ../../bin/yamlmerge .

