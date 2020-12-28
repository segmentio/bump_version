SHELL = /bin/bash -o pipefail
.PHONY: test

export PATH := $(PATH):/usr/local/meter/bin

install:
	go install ./...

lint:
	go vet ./...
	staticcheck ./...

test: lint
	go test -race ./... -timeout 1s

release: install test
	bump_version minor lib/lib.go

ci-install:
	curl -s https://packagecloud.io/install/repositories/meter/public/script.deb.sh | sudo bash
	sudo apt-get -qq -o=Dpkg::Use-Pty=0 install staticcheck

ci: ci-install lint test
