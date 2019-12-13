GO11MODULES=on
APP?=lighter
RELEASE="0.2.5"
COMMIT=$(shell git rev-parse --short HEAD)

all: test

.PHONY: mod
mod:
	go mod tidy
	go mod vendor

.PHONY: test
test:
	go test -v -count=1 -race ./...

.PHONY: tag
tag:
	git tag "v${RELEASE}"
	git push origin "v${RELEASE}"
	git log --oneline
