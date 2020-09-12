#!/bin/sh

set -eu

# OS X
export GOOS=darwin
export GOARCH=amd64

go build -o bin/splable-agent-$VERSION-$GOOS-$GOARCH

# Raspberry Pi
export GOOS=linux
export GOARCH=arm
export GOARM=7

go build -o bin/splable-agent-$VERSION-$GOOS-$GOARCH$GOARM
