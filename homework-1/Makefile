SHELL := /bin/bash

# Пути и утилиты
BIN_DIR := bin
RES_DIR := result
MKDIR := mkdir -p $(BIN_DIR)
RM := rm -f
RMDIR := rm -rf

TOOLS_DIR := $(shell go env GOPATH)/bin
PROTOC_GEN_GO := $(BIN_DIR)/protoc-gen-go
PROTOC_GEN_GO_GRPC := $(BIN_DIR)/protoc-gen-go-grpc
PROTOC_GEN_VALIDATE := $(BIN_DIR)/protoc-gen-validate
PROTOC_GEN_GRPC_GATEWAY := $(BIN_DIR)/protoc-gen-grpc-gateway
PROTOC_GEN_OPENAPIV2 := $(BIN_DIR)/protoc-gen-openapiv2
SWAGGER_CLI := $(BIN_DIR)/swagger

PROTO_DEPS := \
	$(PROTOC_GEN_GO) \
	$(PROTOC_GEN_GO_GRPC) \
	$(PROTOC_GEN_VALIDATE) \
	$(PROTOC_GEN_GRPC_GATEWAY) \
	$(PROTOC_GEN_OPENAPIV2) \
	$(SWAGGER_CLI)

GOCYCLO := $(TOOLS_DIR)/gocyclo
GOCOGNIT := $(TOOLS_DIR)/gocognit

SWAGGER_FILE ?= $(OUT_DIR)/apidocs.swagger.json
SWAGGER_PORT ?= 8081

OUT_DIR := internal/proto/loms
PROTO_DIR := internal/proto

.PHONY: run-all test-coverage install-tools cyclo cognit report clean proto-deps vendor-proto generate-proto check-proto-files link-tools generate-swagger serve-swagger

run-all:
	go run ./cmd/cart/main.go

test-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report saved to coverage.html"

install-tools:
	@if [ ! -f "$(GOCYCLO)" ]; then \
		echo "Installing gocyclo..."; \
		go install github.com/fzipp/gocyclo/cmd/gocyclo@latest; \
	fi
	@if [ ! -f "$(GOCOGNIT)" ]; then \
		echo "Installing gocognit..."; \
		go install github.com/uudashr/gocognit/cmd/gocognit@latest; \
	fi

cyclo: install-tools
	@mkdir -p result
	-@"$(GOCYCLO)" -over 10 . > $(RES_DIR)/cyclo.txt
	@echo "Cyclo report saved to $(RES_DIR)/cyclo.txt"

cognit: install-tools
	@mkdir -p result
	-@"$(GOCOGNIT)" -over 10 . > $(RES_DIR)/cognit.txt
	@echo "Cognit report saved to $(RES_DIR)/cognit.txt"

report: cyclo cognit
	@echo "Reports generated in $(RES_DIR)/"

clean:
	$(RM) coverage.out coverage.html
	$(RMDIR) $(RES_DIR)

vendor-proto:
	@echo "Setting up vendor-proto directories..."
	@mkdir -p vendor-proto

	@if [ ! -d "vendor-proto/google/api" ]; then \
		echo "Cloning googleapis/google/api..."; \
		git clone -b master --single-branch -n --depth=1 --filter=tree:0 \
			https://github.com/googleapis/googleapis vendor-proto/googleapis-tmp && \
		cd vendor-proto/googleapis-tmp && \
		git sparse-checkout set google/api && \
		git checkout && \
		mkdir -p ../google && \
		mv google/api ../google/ && \
		cd ../.. && \
		rm -rf vendor-proto/googleapis-tmp; \
	fi

	@if [ ! -d "vendor-proto/google/protobuf" ]; then \
		echo "Cloning protobuf/google/protobuf..."; \
		git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
			https://github.com/protocolbuffers/protobuf vendor-proto/protobuf-tmp && \
		cd vendor-proto/protobuf-tmp && \
		git sparse-checkout set src/google/protobuf && \
		git checkout && \
		mkdir -p ../google && \
		mv src/google/protobuf ../google/ && \
		cd ../.. && \
		rm -rf vendor-proto/protobuf-tmp; \
	fi

	@if [ ! -d "vendor-proto/google.golang.org/protobuf" ]; then \
		echo "Cloning protobuf-go/types/known..."; \
		git clone -b master --single-branch -n --depth=1 --filter=tree:0 \
			https://github.com/protocolbuffers/protobuf-go vendor-proto/protobuf-go-tmp && \
		cd vendor-proto/protobuf-go-tmp && \
		git sparse-checkout set types/known && \
		git checkout && \
		mkdir -p ../google.golang.org/protobuf && \
		mv types/known ../google.golang.org/protobuf/ && \
		cd ../.. && \
		rm -rf vendor-proto/protobuf-go-tmp; \
	fi

	@if [ ! -d "vendor-proto/google.golang.org/grpc" ]; then \
		echo "Cloning grpc-go cmd/protoc-gen-go-grpc..."; \
		git clone -b master --single-branch -n --depth=1 --filter=tree:0 \
			https://github.com/grpc/grpc-go vendor-proto/grpc-go-tmp && \
		cd vendor-proto/grpc-go-tmp && \
		git sparse-checkout set cmd/protoc-gen-go-grpc && \
		git checkout && \
		mkdir -p ../google.golang.org && \
		mv cmd/protoc-gen-go-grpc ../google.golang.org/grpc && \
		cd ../.. && \
		rm -rf vendor-proto/grpc-go-tmp; \
	fi

	@if [ ! -d "vendor-proto/validate" ]; then \
		echo "Cloning protoc-gen-validate..."; \
		git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
			https://github.com/envoyproxy/protoc-gen-validate vendor-proto/pgv-tmp && \
		cd vendor-proto/pgv-tmp && \
		git sparse-checkout set validate && \
		git checkout && \
		mv validate ../ && \
		cd ../.. && \
		rm -rf vendor-proto/pgv-tmp; \
	fi

	@if [ ! -d "vendor-proto/protoc-gen-openapiv2/options" ]; then \
		echo "Cloning grpc-ecosystem/grpc-gateway (for openapiv2 options)..."; \
		git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
			https://github.com/grpc-ecosystem/grpc-gateway vendor-proto/grpc-gateway-tmp && \
		cd vendor-proto/grpc-gateway-tmp && \
		git sparse-checkout set protoc-gen-openapiv2/options && \
		git checkout && \
		mkdir -p ../protoc-gen-openapiv2 && \
		mv protoc-gen-openapiv2/options ../protoc-gen-openapiv2/ && \
		cd ../.. && \
		rm -rf vendor-proto/grpc-gateway-tmp; \
	fi

	@echo "Vendor-proto setup complete!"

proto-deps: $(PROTO_DEPS)
	@echo "All proto dependencies installed to $(BIN_DIR)/"
	@ls -la $(BIN_DIR)/

$(PROTOC_GEN_GO):
	@echo "Installing protoc-gen-go..."
	@mkdir -p $(BIN_DIR)
	@GOBIN=$(PWD)/$(BIN_DIR) go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

$(PROTOC_GEN_GO_GRPC):
	@echo "Installing protoc-gen-go-grpc..."
	@mkdir -p $(BIN_DIR)
	@GOBIN=$(PWD)/$(BIN_DIR) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

$(PROTOC_GEN_VALIDATE):
	@echo "Installing protoc-gen-validate..."
	@mkdir -p $(BIN_DIR)
	@GOBIN=$(PWD)/$(BIN_DIR) go install github.com/envoyproxy/protoc-gen-validate@latest

$(PROTOC_GEN_GRPC_GATEWAY):
	@echo "Installing protoc-gen-grpc-gateway..."
	@mkdir -p $(BIN_DIR)
	@GOBIN=$(PWD)/$(BIN_DIR) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest

$(PROTOC_GEN_OPENAPIV2):
	@echo "Installing protoc-gen-openapiv2..."
	@mkdir -p $(BIN_DIR)
	@GOBIN=$(PWD)/$(BIN_DIR) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest

$(SWAGGER_CLI):
	@echo "Installing swagger..."
	@mkdir -p $(BIN_DIR)
	@GOBIN=$(PWD)/$(BIN_DIR) go install github.com/go-swagger/go-swagger/cmd/swagger@latest

check-proto-files:
	@if [ ! -f "$(PROTO_DIR)/validate.proto" ]; then \
		echo "Error: validate.proto not found in $(PROTO_DIR)"; \
		echo "Please download it manually and place it in $(PROTO_DIR):"; \
		echo "curl -o $(PROTO_DIR)/validate.proto https://raw.githubusercontent.com/envoyproxy/protoc-gen-validate/main/validate/validate.proto"; \
		exit 1; \
	fi

generate-proto: vendor-proto proto-deps check-proto-files
	@echo "Generating protobuf files..."
	protoc \
		-I "$(PROTO_DIR)" \
		-I "vendor-proto" \
		--plugin=protoc-gen-go="$(PROTOC_GEN_GO)" \
		--go_out="$(OUT_DIR)" \
		--go_opt paths=source_relative \
		--plugin=protoc-gen-go-grpc="$(PROTOC_GEN_GO_GRPC)" \
		--go-grpc_out="$(OUT_DIR)" \
		--go-grpc_opt paths=source_relative \
		--plugin=protoc-gen-validate="$(PROTOC_GEN_VALIDATE)" \
		--validate_out="lang=go,paths=source_relative:$(OUT_DIR)" \
		--plugin=protoc-gen-grpc-gateway="$(PROTOC_GEN_GRPC_GATEWAY)" \
		--grpc-gateway_out="$(OUT_DIR)" \
		--grpc-gateway_opt paths=source_relative \
		--plugin=protoc-gen-openapiv2="$(PROTOC_GEN_OPENAPIV2)" \
		--openapiv2_out="$(OUT_DIR)" \
		--openapiv2_opt allow_merge=true \
		"$(PROTO_DIR)/loms.proto"
	@echo "Proto files generated successfully in $(OUT_DIR)!"
	go mod tidy

serve-swagger: $(SWAGGER_CLI)
	@$(SWAGGER_CLI) serve --no-open --flavor=swagger --port=$(SWAGGER_PORT) $(SWAGGER_FILE)

