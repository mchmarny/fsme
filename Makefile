RELEASE_VERSION="0.2.1"

all: test

mod:
	go mod tidy
	go mod vendor

test:
	go test ./... -v

tag:
	git tag "release-v${RELEASE_VERSION}"
	git push origin "release-v${RELEASE_VERSION}"
	git log --oneline
