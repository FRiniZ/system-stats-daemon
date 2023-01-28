BIN_NAME=ssd-app
BIN_DIR=./bin
DOCKER_IMG=daemon:develop

APP_VERSION=$(shell scripts/version.sh)
GIT_HASH=$(shell git log --format="%h" -n 1)
LDFLAGS=-X main.release=$(APP_VERSION) -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

BUILD_DIR=build

build-linux:
	@echo "Building Linux binary..."
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BIN_NAME)-linux-amd64 ./cmd/ssd

build-osx:
	@echo "Building Linux binary..."
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BIN_NAME)-darwin-amd64 ./cmd/ssd

build-all: build-linux build-osx

build:
	go build -v -o $(BIN_DIR)/$(BIN_NAME) -ldflags "$(LDFLAGS)" ./cmd/daemon

package-linux:
	@echo "Packaging Linux binary..."
	tar -C $(BUILD_DIR) -zcf $(BUILD_DIR)/$(BIN_NAME)-$(APP_VERSION)-linux-amd64.tar.gz $(BIN_NAME)-linux-amd64

package-osx:
	@echo "Packaging OSX binary..."
	tar -C $(BUILD_DIR) -zcf $(BUILD_DIR)/$(BIN_NAME)-$(APP_VERSION)-darwin-amd64.tar.gz $(BIN_NAME)-darwin-amd64

run: build
	$(BIN) -config ./configs/config.toml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race ./internal/... ./pkg/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build run build-img run-img version test lint