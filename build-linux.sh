#!/usr/bin/env bash
basedir=`pwd`
cd cmd/configgen && GOOS=linux GOARCH=amd64 go build -a -ldflags '-extldflags "-static"' -v -o ${basedir}/bin/linux/configgen
exit $?