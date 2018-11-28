#!/usr/bin/env bash

pushd $GOPATH/src/github.com/tekliner/postgres/hack/gendocs
go run main.go
popd
