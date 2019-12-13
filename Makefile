GO11MODULES=on
APP?=firestoredal
RELEASE="0.2.1"
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
	git tag "release-v${RELEASE}-${COMMIT}"
	git push origin "release-v${RELEASE}-${COMMIT}"
	git log --oneline
