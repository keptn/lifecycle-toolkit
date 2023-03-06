ROOT_DIR := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))
VOLUMES := -v $(ROOT_DIR):/src 
# renovate: datasource=docker depName=klakegg/hugo
HUGO_VERSION := 0.105.0-ext
IMAGE := klakegg/hugo:$(HUGO_VERSION)
PORT := 1314
DOCKER_CMD := docker run --rm -t -e HUGO_CACHEDIR=/src/tmp/.hugo

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
#
# renovate: datasource=docker depName=ghcr.io/igorshubovych/markdownlint-cli
MDL_DOCKER_VERSION := v0.33.0
MDL_CMD := docker run -v $(ROOT_DIR)../:/workdir --rm  ghcr.io/igorshubovych/markdownlint-cli:$(MDL_DOCKER_VERSION)  "**/*.md" 

.PHONY: markdownlint markdownlint-fix
markdownlint:
	$(MDL_CMD)

markdownlint-fix:
	$(MDL_CMD) --fix
