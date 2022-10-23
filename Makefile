GO_VERSION_SHORT:=$(shell echo `go version` | sed -E 's/.* go(.*) .*/\1/g')
ifneq ("1.19.1","$(shell printf "$(GO_VERSION_SHORT)\n1.19.1" | sort -V | head -1)")
$(error NEED GO VERSION >= 1.19.1. Found: $(GO_VERSION_SHORT))
endif

export GO111MODULE=on

SERVICE_NAME=act-device-api
SERVICE_PATH=ozonmp/act-device-api

PGV_VERSION:="v0.6.1"
BUF_VERSION:="v0.56.0"

OS_NAME=$(shell uname -s)
OS_ARCH=$(shell uname -m)
GO_BIN=$(shell go env GOPATH)/bin
BUF_EXE=$(GO_BIN)/buf$(shell go env GOEXE)

ifeq ("NT", "$(findstring NT,$(OS_NAME))")
OS_NAME=Windows
endif

.PHONY: run
run:
	go run cmd/grpc-server/main.go

.PHONY: lint
lint:
	golangci-lint run --new-from-rev 09cb7872e296d9e9fa8d4f816bb73cf25706b44a ./...


.PHONY: test
test:
	go test -v -race -timeout 30s -coverprofile cover.out ./...
	go tool cover -func cover.out

.PHONY: gotestsum
gotestsum:
	gotestsum --junitfile unit-tests.xml \
		--jsonfile json-report.txt \
		-- -v -race -timeout 30s -coverprofile cover.out ./...
	go tool cover -func cover.out

# ----------------------------------------------------------------

.PHONY: generate-proto-desc
generate-proto-desc:
	cd api && \
		protoc -o kafka/outfile.desc --include_imports ozonmp/act_device_api/v1/act_device_api.proto

.PHONY: generate
generate: .generate-install-buf .generate-go

.PHONY: generate-go
generate-go: .generate-install-buf .generate-go

.generate-install-buf:
	@ command -v buf 2>&1 > /dev/null || (echo "Install buf" && \
    		mkdir -p "$(GO_BIN)" && \
    		curl -sSL0 https://github.com/bufbuild/buf/releases/download/$(BUF_VERSION)/buf-$(OS_NAME)-$(OS_ARCH)$(shell go env GOEXE) -o "$(BUF_EXE)" && \
    		chmod +x "$(BUF_EXE)")

.generate-go:
	$(BUF_EXE) generate

# ----------------------------------------------------------------

.PHONY: deps
deps: deps-go

.PHONY: deps-go
deps-go:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.27.1
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1.0
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.5.0
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.5.0
	go install github.com/envoyproxy/protoc-gen-validate@$(PGV_VERSION)
	go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger@latest

.PHONY: build
build: generate .build

.PHONY: build-go
build-go: generate-go .build

.build:
	go mod download && CGO_ENABLED=0  go build \
		-tags='no_mysql no_sqlite3' \
		-ldflags=" \
			-X 'gitlab.ozon.dev/$(SERVICE_PATH)/internal/config.version=$(VERSION)' \
			-X 'gitlab.ozon.dev/$(SERVICE_PATH)/internal/config.commitHash=$(COMMIT_HASH)' \
		" \
		-o ./bin/grpc-server$(shell go env GOEXE) ./cmd/grpc-server/main.go


# ----------------------------------------------------------------
.PHONY: docker-build
docker-build:
	docker-compose build act-device-api

.PHONY: dc-serv-up
dc-serv-up:
	docker-compose -f docker-compose.service.yaml up -d

.PHONY: dc-serv-down
dc-serv-down:
	docker-compose -f docker-compose.service.yaml down -v

.PHONY: dc-serv-env-up
dc-serv-env-up:
	docker-compose -f docker-compose.service-env.yaml up -d

.PHONY: dc-serv-env-down
dc-serv-env-down:
	docker-compose -f docker-compose.service-env.yaml down -v -t0

.PHONY: dc-serv-stop
dc-serv-stop:
	docker-compose -f docker-compose.service.yaml stop

.PHONY: dc-serv-env-stop
dc-serv-env-stop:
	docker-compose -f docker-compose.service-env.yaml stop

.PHONY: dc-serv-rebuild-reup
dc-serv-rebuild-reup: dc-serv-down
	docker-compose -f docker-compose.service.yaml up --build --force-recreate -V -d

.PHONY: dc-serv-env-rebuild-reup
dc-serv-env-rebuild-reup: dc-serv-env-down
	docker-compose -f docker-compose.service-env.yaml up --build --force-recreate -V -d

.PHONY: grpc-tests
grpc-http-tests: 
	go test ./tests -v
	
.PHONY: tools-version
tools-version:
	@ curl --version
	@ golangci-lint --version
	@ protoc --version
	@ docker --version
	@ docker-compose --version
