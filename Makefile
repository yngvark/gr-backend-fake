SHELL         = bash
IMAGE = ghcr.io/yngvark/gr-zombie
DOCKER_COMPOSE_FILE=docker-compose-kafka.yaml

GO := $(shell command -v go 2> /dev/null)
ifndef GO
$(error go is required, please install)
endif

# Dependencies
GOPATH			:= $(shell go env GOPATH)
GOBIN			?= $(GOPATH)/bin
GOFUMPT			:= $(GOBIN)/gofumpt
GOLANGCILINT   	:= $(GOBIN)/golangci-lint

# Paths
FILES = $(shell find . -name '.?*' -prune -o -name vendor -prune -o -name '*.go' -print)

PKGS  = $(or $(PKG),$(shell env GO111MODULE=on $(GO) list ./...))
TESTPKGS = $(shell env GO111MODULE=on $(GO) list -f \
            '{{ if or .TestGoFiles .XTestGoFiles }}{{ .ImportPath }}{{ end }}' \
            $(PKGS))

# Directories
BUILD_DIR     := build

.PHONY: help
help: ## Print this menu
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

init: ## - Set up stuff you need to run locally
	@echo "Install mkcert and run:"
	@echo "mkcert localhost"

gofumpt: ## -
	$(GO) get -u mvdan.cc/gofumpt

fmt: gofumpt  ## -
	$(GO) fmt $(PKGS)
	$(GOFUMPT) -s -w $(FILES)

golangcilint:
	# To bump, simply change the version at the end to the desired version. The git sha here points to the newest commit
	# of the install script verified by our team located here: https://github.com/golangci/golangci-lint/blob/master/install.sh
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/17d24ebd671875cdf52804e1ca72ca8f0718a844/install.sh | sh -s -- -b ${GOBIN} v1.42.1

lint: golangcilint ## -
	$(GOLANGCILINT) run

check: fmt lint test

test:
	go test $(TESTPKGS)

build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)

build-docker: ## -
	docker build . -t $(IMAGE)

run: ##-
#	PORT="8080" \
#	LOG_TYPE="simple" \
#	GAME_QUEUE_TYPE="websocket" \

	ALLOWED_CORS_ORIGINS="http://localhost:3000,http://localhost:3001,http://localhost:30010,http://127.0.0.1:3000" \
	go run *.go

run-docker: build-docker ## -
	docker run \
		 --name zombie-backend \
		 --rm \
		 -e PORT="8080" \
		 -e LOG_TYPE="simple" \
		 -e ALLOWED_CORS_ORIGINS="http://localhost:3001,http://localhost:30010" \
		 -p 8080:8080 \
		 $(IMAGE)

push: build-docker ## -
	@echo Remember to login with docker login ghcr.io. Password is personal access token made in Github.
	@echo
	docker push $(IMAGE)

up: ## docker-compose up -d with logs
	#(docker-compose -f ${DOCKER_COMPOSE_FILE} down || true) && \
	docker-compose -f ${DOCKER_COMPOSE_FILE} up -d && \
	docker logs -f gr-zombie_broker_1

down: ## docker-compose down
	docker-compose -f ${DOCKER_COMPOSE_FILE} down

ws-edit: ## - Edit websocket config
	docker run -it -v zombie-go_pulsarconf:/pconf yngvark/linuxtools vim /pconf/websocket.conf

# Coverage
GOCOVMERGE      := $(GOBIN)/gocovmerge
GOCOVXML        := $(GOBIN)/gocov-xml
GOCOV           := $(GOBIN)/gocov

COVERAGE_MODE    = atomic
COVERAGE_PROFILE = $(COVERAGE_DIR)/profile.out
COVERAGE_XML     = $(COVERAGE_DIR)/coverage.xml
COVERAGE_HTML    = $(COVERAGE_DIR)/index.html

$(GOCOVMERGE):
	$(GO) install github.com/wadey/gocovmerge@latest

$(GOCOVXML):
	$(GO) install github.com/AlekSi/gocov-xml@latest

$(GOCOV):
	$(GO) install github.com/axw/gocov/gocov@v1.0.0

test-coverage-tools: | $(GOCOVMERGE) $(GOCOV) $(GOCOVXML)
test-coverage: COVERAGE_DIR := $(BUILD_DIR)/test/coverage.$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
test-coverage: fmt lint test-coverage-tools
	@mkdir -p $(COVERAGE_DIR)/coverage
	@for pkg in $(TESTPKGS); do \
        go test \
            -coverpkg=$$(go list -f '{{ join .Deps "\n" }}' $$pkg | \
                    grep '^$(MODULE)/' | \
                    tr '\n' ',')$$pkg \
            -covermode=$(COVERAGE_MODE) \
            -coverprofile="$(COVERAGE_DIR)/coverage/`echo $$pkg | tr "/" "-"`.cover" $$pkg ;\
     done
	@$(GOCOVMERGE) $(COVERAGE_DIR)/coverage/*.cover > $(COVERAGE_PROFILE)
	@$(GO) tool cover -html=$(COVERAGE_PROFILE) -o $(COVERAGE_HTML)
	@$(GOCOV) convert $(COVERAGE_PROFILE) | $(GOCOVXML) > $(COVERAGE_XML)
