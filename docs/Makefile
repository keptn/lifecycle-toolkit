DOCREPO := github.com/keptn-sandbox/lifecycle-toolkit-docs
TMPDIR := $(CURDIR)/tmp
VOLUMES := -v $(TMPDIR)/lifecycle-toolkit-docs:/src -v $(TMPDIR)/lifecycle-toolkit-tmp/public:/target -v $(CURDIR)/en:/src/content/en
# renovate: datasource=docker depName=klakegg/hugo
IMAGE := klakegg/hugo:0.105.0-ext

.PHONY: build server clean htmltest

clone:
	@rm -rf $(TMPDIR)/lifecycle-toolkit-docs | true
	@mkdir -p $(TMPDIR)/lifecycle-toolkit-docs
	@git clone https://$(DOCREPO) $(TMPDIR)/lifecycle-toolkit-docs

build:
	docker run --rm -it $(VOLUMES) $(IMAGE) -D -v

server:
	docker run --rm -it $(VOLUMES) -p 1313:1313 $(IMAGE) server -D

clean:
	docker run --rm -it $(VOLUMES) $(IMAGE) --cleanDestinationDir

htmltest:
	docker run -v $(CURDIR):/test --rm wjdp/htmltest -s -c .htmltest.yml public
