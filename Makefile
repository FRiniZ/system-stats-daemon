BINARY_SSD=ssdaemon-app
BINARY_SSC=ssclient-app
BUILD_DIR=build

APP_VERSION=$(shell scripts/version.sh)
LDFLAGS=-X main.release=$(APP_VERSION) -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

GO_BUILD_CMD= CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)"

.PHONY: all
all: clean lint test build-all package-all

.PHONY: lint
install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.50.1

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: test
test:
	@echo "Running tests..."
	@go test ./...

.PHONY: pre-build
pre-build:
	@mkdir -p $(BUILD_DIR)

.PHONY: build-linux
build-linux: pre-build
	@echo "Building Linux binary..."
	GOOS=linux GOARCH=amd64 $(GO_BUILD_CMD) -o $(BUILD_DIR)/$(BINARY_SSD)-linux-amd64 ./cmd/ssdaemon
	GOOS=linux GOARCH=amd64 $(GO_BUILD_CMD) -o $(BUILD_DIR)/$(BINARY_SSC)-linux-amd64 ./cmd/ssclient

.PHONY: build-osx
build-osx: pre-build
	@echo "Building OSX binary..."
	GOOS=darwin GOARCH=amd64 $(GO_BUILD_CMD) -o $(BUILD_DIR)/$(BINARY_SSD)-darwin-amd64 ./cmd/ssdaemon
	GOOS=darwin GOARCH=amd64 $(GO_BUILD_CMD) -o $(BUILD_DIR)/$(BINARY_SSC)-darwin-amd64 ./cmd/ssclient

.PHONY: build build-all
build-all: build-linux build-osx

.PHONY: package-linux
package-linux:
	@echo "Packaging Linux binary..."
	tar -C $(BUILD_DIR) -zcf $(BUILD_DIR)/$(BINARY_SSD)-$(APP_VERSION)-linux-amd64.tar.gz $(BINARY_SSD)-linux-amd64
	tar -C $(BUILD_DIR) -zcf $(BUILD_DIR)/$(BINARY_SSC)-$(APP_VERSION)-linux-amd64.tar.gz $(BINARY_SSC)-linux-amd64

.PHONY: package-osx
package-osx:
	@echo "Packaging OSX binary..."
	tar -C $(BUILD_DIR) -zcf $(BUILD_DIR)/$(BINARY_SSD)-$(APP_VERSION)-darwin-amd64.tar.gz $(BINARY_SSD)-darwin-amd64
	tar -C $(BUILD_DIR) -zcf $(BUILD_DIR)/$(BINARY_SSC)-$(APP_VERSION)-darwin-amd64.tar.gz $(BINARY_SSC)-darwin-amd64

.PHONY: package-all
package-all: package-linux package-osx

.PHONY: docker
docker:
	docker build --force-rm -t $(BINARY_SSD) .

.PHONY: build-in-docker
build-in-docker: docker
	docker rm -f $(BINARY_SSD) || true
	docker create --name $(BINARY_SSD) $(BINARY_SSD)
	docker cp '$(BINARY_SSD):/opt/' $(BUILD_DIR)
	docker rm -f $(BINARY_SSD)

.PHONY: clean
clean:
	@echo "Cleaning..."
	@rm -Rf $(BUILD_DIR)


.PHONY: build
build:
	@echo "Building native binary..."
	$(GO_BUILD_CMD) -o $(BUILD_DIR)/$(BINARY_SSD)-native-amd64 ./cmd/ssdaemon
	$(GO_BUILD_CMD) -o $(BUILD_DIR)/$(BINARY_SSC)-native-amd64 ./cmd/ssclient


.PHONY: generate
generate:
	go generate ./...