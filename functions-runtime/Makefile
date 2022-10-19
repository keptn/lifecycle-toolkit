# RELEASE_REGISTRY is the container registry to push
# into.
RELEASE_REGISTRY?=ghcr.io/keptn-sandbox
RELEASE_VERSION?=$(shell date +%Y%m%d%s)-v0.24.3#$(shell git describe --tags --match "v*")
TAG?=latest
RELEASE_IMAGE:=functions-runtime:$(TAG)

ARCHS = amd64 arm64
COMMONENVVAR=GOOS=$(shell uname -s | tr A-Z a-z)
BUILDENVVAR=CGO_ENABLED=0

# The RELEASE_VERSION variable can have one of two formats:
# v20201009-v0.18.800-46-g939c1c0 - automated build for a commit(not a tag) and also a local build
# v20200521-v0.18.800             - automated build for a tag
VERSION=$(shell echo $(RELEASE_VERSION) | awk -F - '{print $$2}')

.PHONY: clean
clean:
	rm -rf ./bin
