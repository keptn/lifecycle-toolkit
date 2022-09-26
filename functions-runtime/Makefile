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


.PHONY: release-image
release-image: release-image.amd64 release-image.arm64v8

.PHONY: release-image.amd64
release-image.amd64: clean
	docker build --no-cache --build-arg ARCH="amd64" -t $(RELEASE_REGISTRY)/$(RELEASE_IMAGE)-amd64 .

.PHONY: release-image.arm64v8
release-image.arm64v8: clean
	docker build  --no-cache --build-arg ARCH="arm64v8" -t $(RELEASE_REGISTRY)/$(RELEASE_IMAGE)-arm64 .

.PHONY: push-release-images
push-release-images: release-image.amd64 release-image.arm64v8
	#gcloud auth configure-docker
	for arch in $(ARCHS); do \
		docker push $(RELEASE_REGISTRY)/$(RELEASE_IMAGE)-$${arch} ;\
	done
	DOCKER_CLI_EXPERIMENTAL=enabled docker manifest create $(RELEASE_REGISTRY)/$(RELEASE_IMAGE) $(addprefix --amend $(RELEASE_REGISTRY)/$(RELEASE_IMAGE)-, $(ARCHS))
	for arch in $(ARCHS); do \
		DOCKER_CLI_EXPERIMENTAL=enabled docker manifest annotate --arch $${arch} $(RELEASE_REGISTRY)/$(RELEASE_IMAGE) $(RELEASE_REGISTRY)/$(RELEASE_IMAGE)-$${arch} ;\
	done
	DOCKER_CLI_EXPERIMENTAL=enabled docker manifest push $(RELEASE_REGISTRY)/$(RELEASE_IMAGE) ;\

.PHONY: clean
clean:
	rm -rf ./bin