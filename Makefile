include $(DEVHOME)/scripts/Makefile

SHELL=/bin/bash
GIT_COMMIT_SHORT = $(shell git rev-parse --short HEAD)

prettier:
	yarn run format

code:
	code ~/workspace/config/vs-code/engineering.code-workspace

install:
	export CGO_ENABLED=0 && \
	go install && \
	go get

run:
	go run .

clean:
	rm -rf bin/*

build: clean
	export VERSION=$(GIT_COMMIT_SHORT) && \
	.buildkite/build.sh

fmt:
	go fmt

scp:
	scp bin/splable-agent-$(GIT_COMMIT_SHORT)-linux-arm7 pi@192.168.1.60:/home/pi/splable-agent
