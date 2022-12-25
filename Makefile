RELEASE_VERSION ?=$(shell cat ./.version)

all: help

version: ## Prints the current version
	@echo $(RELEASE_VERSION)
.PHONY: version

tidy: ## Updates the go modules and vendors all dependancies 
	go mod tidy
	go mod vendor
.PHONY: tidy

upgrade: ## Upgrades all dependancies 
	go get -d -u ./...
	go mod tidy
	go mod vendor
.PHONY: upgrade

test: tidy ## Runs unit tests
		go test -count=1 -race -covermode=atomic -coverprofile=cover.out ./...
.PHONY: test

run: tidy ## Runs uncompiled version of the app
	go run cmd/disco/main.go
.PHONY: run

cover: test ## Runs unit tests and putputs coverage
	go tool cover -func=cover.out
.PHONY: cover

lint: ## Lints the entire project 
	golangci-lint -c .golangci.yaml run
.PHONY: lint

release: tidy ## Builds CLI binary
	goreleaser release --snapshot --rm-dist --timeout 10m0s
	mkdir -p ./bin
	mv dist/disco-arm64-darwin ./bin/disco
.PHONY: release

img-run: ## Runs run image discovery command on the compiled version of CLI
	debug=true bin/disco run img -o test.json
.PHONY: img-run

vul-run: ## Runs image vulnerability discovery command on the compiled version of CLI
	debug=true bin/disco run vul -o test.json
.PHONY: vul-run

vul-local-run: ## Runs image vulnerability discovery command on the compiled version of CLI
	debug=true bin/disco run vul --local -o test.json
.PHONY: vul-run

cve-run: ## Runs image vulnerability discovery command on the compiled version of CLI
	debug=true bin/disco run vul --cve CVE-2017-11164
.PHONY: cve-run

cve-local-run: ## Runs image vulnerability discovery command on the compiled version of CLI
	debug=true bin/disco run vul --local --cve CVE-2017-11164
.PHONY: cve-run

lic-run: ## Runs run license discovery command on the compiled version of CLI
	debug=true bin/disco run lic -o test.json
.PHONY: lic-run

tag: ## Creates release tag 
	git tag -s -m "release $(RELEASE_VERSION)" $(RELEASE_VERSION)
	git push origin $(RELEASE_VERSION)
.PHONY: tag

tagless: ## Delete the current release tag 
	git tag -d $(RELEASE_VERSION)
	git push --delete origin $(RELEASE_VERSION)
.PHONY: tagless

clean: ## Cleans bin and temp directories
	go clean
	rm -fr ./vendor
	rm -fr ./bin
.PHONY: clean

help: ## Display available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk \
		'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
.PHONY: help