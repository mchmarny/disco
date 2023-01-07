RELEASE_VERSION :=$(shell cat .version)
COMMIT          :=$(shell git rev-parse HEAD)
YAML_FILES      :=$(shell find . -type f -regex ".*yaml" -print)
CURRENT_DATE	:=$(shell date '+%Y-%m-%dT%H:%M:%SZ')

## Variable assertions
ifndef RELEASE_VERSION
	$(error RELEASE_VERSION is not set)
endif

ifndef COMMIT
	$(error COMMIT is not set)
endif

all: help

.PHONY: version
version: ## Prints the current version
	@echo $(RELEASE_VERSION)

.PHONY: tidy
tidy: ## Updates the go modules and vendors all dependancies 
	go mod tidy
	go mod vendor

.PHONY: upgrade
upgrade: ## Upgrades all dependancies 
	go get -d -u ./...
	go mod tidy
	go mod vendor

.PHONY: test
test: tidy ## Runs unit tests
	mkdir -p tmp
	go test -short -count=1 -race -covermode=atomic -coverprofile=cover.out ./...

.PHONY: cover
cover: test ## Runs unit tests and putputs coverage
	go tool cover -func=cover.out

.PHONY: vulncheck
vulncheck: test ## Runs go vulneerability check locally
	govulncheck ./...

.PHONY: lint
lint: lint-go lint-yaml ## Lints the entire project 
	@echo "Completed Go and YAML lints"

.PHONY: lint
lint-go: ## Lints the entire project using go 
	golangci-lint -c .golangci.yaml run

.PHONY: lint-yaml
lint-yaml: ## Runs yamllint on all yaml files (brew install yamllint)
	yamllint -c .yamllint $(YAML_FILES)

.PHONY: build
build: build-server build-cli ## Builds binaries
	@echo "Completed release"

.PHONY: build-server
build-server: tidy ## Builds Server binary
	mkdir -p ./bin
	CGO_ENABLED=0 go build -trimpath \
    -ldflags="-w -s -X main.version=$(RELEASE_VERSION) -extldflags '-static'" \
    -a -mod vendor -o bin/server cmd/server/main.go

.PHONY: build-cli
build-cli: tidy ## Builds CLI binary
	mkdir -p ./bin
	CGO_ENABLED=0 go build -trimpath -ldflags="\
    -w -s -X main.version=$(RELEASE_VERSION) \
	-w -s -X main.commit=$(COMMIT) \
	-w -s -X main.date=$(CURRENT_DATE) \
	-extldflags '-static'" \
    -a -mod vendor -o bin/disco cmd/cli/main.go

.PHONY: server
server: ## Runs previsouly built server binary
	tools/server 

.PHONY: image
image: ## Builds new image locally (dirty)
	tools/image

.PHONY: invoke
invoke: ## Invokes the service
	tools/invoke

.PHONY: run
run: ## Runs bash on latest artomator image
	tools/run

.PHONY: release
release: test lint vulncheck tag ## Runs test, lint, vulncheck, and tag before release
	@echo "Releasing: $(RELEASE_VERSION)"
	tools/gh-wait
	tools/tf-apply

.PHONY: infra
infra: ## Applies Terraform
	terraform -chdir=./deploy apply -auto-approve

.PHONY: nice
infra-fmt: ## Formats Terraform
	terraform -chdir=./deploy fmt

.PHONY: tag
tag: ## Creates release tag 
	git tag -s -m "release $(RELEASE_VERSION)" $(RELEASE_VERSION)
	git push origin $(RELEASE_VERSION)

.PHONY: tagless
tagless: ## Delete the current release tag 
	git tag -d $(RELEASE_VERSION)
	git push --delete origin $(RELEASE_VERSION)

.PHONY: clean
clean: ## Cleans bin and temp directories
	go clean
	rm -fr ./vendor
	rm -fr ./bin

.PHONY: docker-clean
docker-clean: ## Removes orpaned docker volumes
	@echo "stopping all containers..."
	docker stop $(shell docker ps -aq)
	@echo "removing all containers..." 
	docker rm $(shell docker ps -aq)
	@echo "prunning system..."
	docker system prune -a --volumes
	@echo "done"


.PHONY: help
help: ## Display available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk \
		'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
