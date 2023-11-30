MKDOCS_DOCKER_IMAGE=squidfunk/mkdocs-material:9
PWD=$(shell pwd)

.PHONY: build
build:
	docker run --rm -it -v ${PWD}:/docs --entrypoint "" ${MKDOCS_DOCKER_IMAGE} sh -c 'pip3 install -r /docs/requirements.txt; mkdocs build'

.PHONY: serve
serve: build
	docker run --rm -it -p 8000:8000 -v ${PWD}/site:/site -w /site --entrypoint "" ${MKDOCS_DOCKER_IMAGE} python3 -m http.server --bind 0.0.0.0 8000
