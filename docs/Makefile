ROOT_DIR := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))
VOLUMES := -v $(ROOT_DIR):/src 
# renovate: datasource=docker depName=klakegg/hugo
HUGO_VERSION := 0.105.0-ext
IMAGE := klakegg/hugo:$(HUGO_VERSION)
PORT := 1314

.PHONY: build server clean shell htmltest

build:
	docker run --rm -it $(VOLUMES) $(IMAGE) -D -v

shell:
	docker run --rm -it $(VOLUMES) $(IMAGE) shell

server:
	docker run --rm -it $(VOLUMES) -p  $(PORT):$(PORT) $(IMAGE) server -D -p $(PORT)

clean:
	docker run --rm -it $(VOLUMES) $(IMAGE) --cleanDestinationDir

htmltest: build
	docker run -v $(CURDIR):/test --rm wjdp/htmltest -s /test/tmp/lifecycle-toolkit-docs/public

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
