# Not Support Windows

.PHONY: help wire gen ent build api openapi init all vendor dep test cover vet lint docker

ifeq ($(OS),Windows_NT)
    IS_WINDOWS := 1
endif

CURRENT_DIR	:= $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))
ROOT_DIR	:= $(dir $(realpath $(lastword $(MAKEFILE_LIST))))
KRATOS_THIRD_PARTY := $(shell go env GOPATH)/pkg/mod/github.com/go-kratos/kratos/v2@v2.9.1/third_party

SRCS_MK		:= $(foreach dir, app, $(wildcard $(dir)/*/*/Makefile))

# use docker compose to run backend services and all its dependency services like redis, mysql, etc.
compose-up:
	@docker compose up -d --force-recreate

compose-down:
	@docker compose down

# download dependencies of module
dep:
	@go mod download

# create vendor
vendor:
	@go mod vendor

# run tests
test:
	@go test ./...

# run coverage tests
cover:
	@go test -v ./... -coverprofile=coverage.out

# run static analysis
vet:
	@go vet

# run lint
lint:
	@golangci-lint run

# generate wire code
wire:
	@go run -mod=mod github.com/google/wire/cmd/wire ./cmd/server

# generate config proto go code
config:
	@protoc -I api/protos --go_out=paths=source_relative:api/gen/go api/protos/conf/recommend.proto

# generate protobuf api go code
api:
	@protoc -I api/third_party -I api/protos -I $(KRATOS_THIRD_PARTY) \
		--go_out=paths=source_relative:api/gen/go \
		--go-grpc_out=paths=source_relative:api/gen/go \
		--go-http_out=paths=source_relative:api/gen/go \
		--grpc-gateway_out=paths=source_relative:api/gen/go \
		--go-mcp_out=paths=source_relative,package_suffix=:api/gen/go \
		api/protos/recommend/service/v1/recommend.proto
	@perl -0pi -e 's/type RecommendServiceServer interface/type RecommendServiceMCPServer interface/g; s/srv RecommendServiceServer/srv RecommendServiceMCPServer/g; s/RegisterRecommendServiceHandler/RegisterRecommendServiceMCPHandler/g; s/type RecommendServiceClient interface/type RecommendServiceMCPClient interface/g; s/client RecommendServiceClient/client RecommendServiceMCPClient/g; s/type ConnectRecommendServiceClient interface/type ConnectRecommendServiceMCPClient interface/g; s/client ConnectRecommendServiceClient/client ConnectRecommendServiceMCPClient/g; s/ForwardToConnectRecommendServiceClient/ForwardToConnectRecommendServiceMCPClient/g; s/ForwardToRecommendServiceClient/ForwardToRecommendServiceMCPClient/g' api/gen/go/recommend/service/v1/recommend.pb.mcp.go

# generate OpenAPI v3 docs.
openapi:
	@cd api && \
	buf generate --template buf.openapi.gen.yaml

# run application
run: config api openapi
	@go run ./cmd/server -conf ./configs

# show help
help:
	@echo ""
	@echo "Usage:"
	@echo " make [target]"
	@echo ""
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
