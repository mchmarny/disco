RELEASE_VERSION ?=$(shell cat .version)
YAML_FILES      :=$(shell find . -type f -regex ".*yaml" -print)

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
	mkdir -p tmp
	go test -count=1 -race -covermode=atomic -coverprofile=cover.out ./...
.PHONY: test

cover: test ## Runs unit tests and putputs coverage
	go tool cover -func=cover.out
.PHONY: cover

lint: lint-go lint-yaml ## Lints the entire project 
	@echo "Completed Go and YAML lints"
.PHONY: lint

lint-go: ## Lints the entire project using go 
	golangci-lint -c .golangci.yaml run
.PHONY: lint

# brew install yamllint
lint-yaml: ## Runs yamllint on all yaml files
	yamllint -c .yamllint $(YAML_FILES)
.PHONY: lint-yaml

build: release-server release-cli ## Builds binaries
	@echo "Completed release"
.PHONY: build

build-server: tidy ## Builds Server binary
	mkdir -p ./bin
	CGO_ENABLED=0 go build -trimpath \
    -ldflags="-w -s -X main.version=$(VERSION) -extldflags '-static'" \
    -a -mod vendor -o ./bin/server cmd/server/main.go
.PHONY: build-server

build-cli: tidy ## Builds CLI binary
	goreleaser release --snapshot --rm-dist --timeout 10m0s
	mkdir -p ./bin
	mv dist/disco_darwin_all/disco ./bin/disco
.PHONY: build-cli

server: ## Runs previsouly built server binary
	tools/server 
.PHONY: server

image: ## Builds new image locally (dirty)
	tools/image
.PHONY: image

run: ## Runs bash on latest artomator image
	tools/run
.PHONY: run

release: test lint tag ## Runs test, lint, and tag before release
	@echo "Releasing: $(VERSION)"
	tools/gh-wait
	tools/tf-apply
.PHONY: release

infra: ## Applies Terraform
	terraform -chdir=./deployment apply -auto-approve
.PHONY: infra

infra-fmt: ## Formats Terraform
	terraform -chdir=./deployment fmt
.PHONY: nice

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

docker-clean: ## Removes orpaned docker volumes
	@echo "stopping all containers..."
	docker stop $(shell docker ps -aq)
	@echo "removing all containers..." 
	docker rm $(shell docker ps -aq)
	@echo "prunning system..."
	docker system prune -a --volumes
	@echo "done"
.PHONY: docker-clean


help: ## Display available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk \
		'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
.PHONY: help