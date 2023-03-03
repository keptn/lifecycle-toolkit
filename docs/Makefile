ROOT_DIR := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))
VOLUMES := -v $(ROOT_DIR):/src 
# renovate: datasource=docker depName=klakegg/hugo
HUGO_VERSION := 0.105.0-ext
IMAGE := klakegg/hugo:$(HUGO_VERSION)
PORT := 1314
DOCKER_CMD := docker run --rm -t -e HUGO_CACHEDIR=/src/tmp/.hugo

.PHONY: build server clean shell htmltest

build:
	$(DOCKER_CMD) $(VOLUMES) $(IMAGE) -D -v

shell:
	$(DOCKER_CMD) -i $(VOLUMES) $(IMAGE) shell

server:
	$(DOCKER_CMD) $(VOLUMES) -p  $(PORT):$(PORT) $(IMAGE) server -D -p $(PORT)

clean:
	$(DOCKER_CMD) $(VOLUMES) $(IMAGE) --cleanDestinationDir

htmltest: build
	$(DOCKER_CMD) -v $(ROOT_DIR):/test wjdp/htmltest -c .htmltest.yml public

.PHONY: lint lint-fix
lint: markdownlint
lint: lint-fix

# Markdown lint configuration
#
# - .markdownlintignore holds the configuration for files to be ignored
# - .markdownlint.yaml contains the rules for markdownfiles
#
MARKDOWNLINT_CMD := docker run -v $(ROOT_DIR)../:/workdir --rm  ghcr.io/igorshubovych/markdownlint-cli:latest  "**/*.md" 

.PHONY: markdownlint markdownlint-fix
markdownlint:
	$(MARKDOWNLINT_CMD)

markdownlint-fix:
	$(MARKDOWNLINT_CMD) --fix
