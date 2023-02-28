ROOT_DIR := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))
VOLUMES := -v $(ROOT_DIR):/src 
# renovate: datasource=docker depName=klakegg/hugo
HUGO_VERSION := 0.105.0-ext
IMAGE := klakegg/hugo:$(HUGO_VERSION)
PORT := 1314

.PHONY: build server clean shell htmltest

build:
	docker run --rm -t $(VOLUMES) $(IMAGE) -D -v

shell:
	docker run --rm -it $(VOLUMES) $(IMAGE) shell

server:
	docker run --rm -t $(VOLUMES) -p  $(PORT):$(PORT) $(IMAGE) server -D -p $(PORT)

clean:
	docker run --rm -t $(VOLUMES) $(IMAGE) --cleanDestinationDir

htmltest: build
	docker run --rm -t -v $(ROOT_DIR):/test wjdp/htmltest -c .htmltest.yml public
