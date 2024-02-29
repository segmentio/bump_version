SHELL = /bin/bash -o pipefail
.PHONY: test

install:
	go install ./...

lint:
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest ./...

test: lint
	go test -race ./... -timeout 1s

release: install test
	next=$$(next_version minor lib/lib.go) && \
		 bump_version minor lib/lib.go && \
		 git tag "v$${next}"

ci: lint test
