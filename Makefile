CUR_DIR         := .
BUILD_PATH      := bin
DOCKERFILE      := $(CUR_DIR)/Dockerfile

PROJECT		:= example.messaging.service
IMAGE_VERSION	:= $(shell git describe --tags --always || echo pre-commit)

SERVER_BINARY   := $(BUILD_PATH)/server
SERVER_PATH     := $(CUR_DIR)/cmd/server

GO_PACKAGES     := $(shell go list ./... | grep -v vendor)

.PHONY: docker
docker:
	@docker build -f $(DOCKERFILE) -t $(PROJECT):$(IMAGE_VERSION) .

.PHONY: build
build:
	@go build -i -v -o $(SERVER_BINARY) $(SERVER_PATH)

.PHONY: clean
clean:
	@rm $(CUR_DIR)/bin/server

.PHONY: test
test:
	@go test -v $(GO_PACKAGES)
