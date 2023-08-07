ROOT_DIR := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))
VOLUMES := -v $(ROOT_DIR):/src 
# renovate: datasource=docker depName=klakegg/hugo
HUGO_VERSION := 0.111.3
IMAGE := klakegg/hugo:$(HUGO_VERSION)-ext
PORT := 1314

INTERACTIVE:=$(shell [ -t 0 ] && echo 1)

ifdef INTERACTIVE
# is a terminal
INTERACTIVE_FLAG := -it
else
INTERACTIVE_FLAG := -t
endif
DOCKER_CMD := docker run --rm $(INTERACTIVE_FLAG) -e HUGO_CACHEDIR=/src/tmp/.hugo -e HUGOxPARAMSxGITHUB_REPO=""

.PHONY: build server clean shell

build:
	$(DOCKER_CMD) $(VOLUMES) $(IMAGE) -D -v

shell:
	$(DOCKER_CMD) -i $(VOLUMES) $(IMAGE) shell

server:
	$(DOCKER_CMD) $(VOLUMES) -p  $(PORT):$(PORT) $(IMAGE) server -D -p $(PORT)

clean:
	$(DOCKER_CMD) $(VOLUMES) $(IMAGE) --cleanDestinationDir

.PHONY: htmltest

# renovate: datasource=docker depName=wjdp/htmltest
HTMLTEST_VERSION := v0.17.0
htmltest: build
	$(DOCKER_CMD) -v $(ROOT_DIR):/test wjdp/htmltest:$(HTMLTEST_VERSION) -c .htmltest.yml public

.PHONY: lint lint-fix
lint: markdownlint
lint: lint-fix

# Markdown lint configuration
#
# - .markdownlintignore holds the configuration for files to be ignored
# - .markdownlint.yaml contains the rules for markdownfiles
# renovate: datasource=docker depName=davidanson/markdownlint-cli2-rules
MDL_DOCKER_VERSION := v0.6.0 # TODO: remove next line on version increase from v0.6.0
MDL_DOCKER_VERSION := next
MDL_CMD := docker run -v $(ROOT_DIR)../:/workdir --rm 

.PHONY: markdownlint markdownlint-fix
markdownlint:
	$(MDL_CMD) davidanson/markdownlint-cli2-rules:${MDL_DOCKER_VERSION} "**/*.md" 

markdownlint-fix:
	$(MDL_CMD) --entrypoint="markdownlint-cli2-fix" davidanson/markdownlint-cli2-rules:${MDL_DOCKER_VERSION} "**/*.md" 
