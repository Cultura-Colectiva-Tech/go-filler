#!/bin/bash

LINUX_BINARY=go-filler-linux
OSX_BINARY=go-filler-osx

if [ -f $LINUX_BINARY ]; then
  rm -v $LINUX_BINARY
fi

if [ -f $OSX_BINARY ]; then
  rm -v $OSX_BINARY
fi

echo "Making Linux 64bit binary: ${LINUX_BINARY}"
GOOS=linux GOARCH=amd64 go build -v -o $LINUX_BINARY

echo "Making OSX 64bit binary: ${OSX_BINARY}"
GOOS=darwin GOARCH=amd64 go build -v -o $OSX_BINARY
