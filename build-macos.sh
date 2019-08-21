#!/usr/bin/env bash
basedir=`pwd`
cd cmd/configgen && GOOS=darwin GOARCH=amd64 go build -v -o ${basedir}/bin/darwin/configgen
exit $?