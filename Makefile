RELEASE_VERSION="0.2.1"

all: test

mod:
	go mod download
	go mod tidy

test:
	go test ./... -v

tag:
	git tag "release-v${RELEASE_VERSION}"
	git push origin "release-v${RELEASE_VERSION}"
	git log --oneline
