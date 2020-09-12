include $(BUILDKITE_BUILD_CHECKOUT_PATH)/ci-scripts/Makefile

AWS_ASSUME_ROLE_NAME = PackerHandcrafted
VERSION = $(shell jq -r .version ../version.json)

PRETTIER = docker run --rm -it \
	-v "$(BUILDKITE_BUILD_CHECKOUT_PATH):/work" \
	tmknom/prettier:2.0.5

GO = docker run --rm -it \
	-v "$(BUILDKITE_BUILD_CHECKOUT_PATH):/src" \
	-w /src \
	-e VERSION=v$(VERSION) \
	golang:1.14.9-alpine3.12

prettier:
	@echo "--- :prettier: Prettier"
	$(PRETTIER) --check .

fmt:
	@echo "--- :go: Fmt"
	$(GO) go fmt

build:
	@echo "--- :go: Build"
	$(GO) /src/.buildkite/build.sh
