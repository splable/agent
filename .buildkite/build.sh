#!/bin/sh

set -eu

# OS X
export CGO_ENABLED=0
export GOOS=darwin
export GOARCH=amd64

go build -o bin/splable-agent-$VERSION-$GOOS-$GOARCH

# Raspberry Pi
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=arm
export GOARM=7

go build -o bin/splable-agent-$VERSION-$GOOS-$GOARCH$GOARM
