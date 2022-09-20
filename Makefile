SHELL := /bin/bash
BASEDIR = $(shell pwd)

# build with version infos
versionDir = "github.com/go-eagle/eagle/pkg/version"
gitTag = $(shell if [ "`git describe --tags --abbrev=0 2>/dev/null`" != "" ];then git describe --tags --abbrev=0; else git log --pretty=format:'%h' -n 1; fi)
buildDate = $(shell TZ=Asia/Shanghai date +%FT%T%z)
gitCommit = $(shell git log --pretty=format:'%H' -n 1)
gitTreeState = $(shell if git status|grep -q 'clean';then echo clean; else echo dirty; fi)

ldflags="-w -X ${versionDir}.gitTag=${gitTag} -X ${versionDir}.buildDate=${buildDate} -X ${versionDir}.gitCommit=${gitCommit} -X ${versionDir}.gitTreeState=${gitTreeState}"

PROJECT_NAME := "github.com/go-microservice/user-service"
PKG := "$(PROJECT_NAME)"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

# proto
APP_RELATIVE_PATH=$(shell a=`basename $$PWD` && echo $$b)
API_PROTO_FILES=$(shell find api$(APP_RELATIVE_PATH) -name *.proto)

# init environment variables
export PATH        := $(shell go env GOPATH)/bin:$(PATH)
export GOPATH      := $(shell go env GOPATH)
export GO111MODULE := on

# make   make all
.PHONY: all
all: lint test build

.PHONY: build
# make build, Build the binary file
build: dep
	@go build -v -ldflags ${ldflags} .

.PHONY: dep
# make dep Get the dependencies
dep:
	@go mod download

.PHONY: fmt
# make fmt
fmt:
	@gofmt -s -w .

.PHONY: lint
# make lint
lint:
	@golint -set_exit_status ${PKG_LIST}

.PHONY: ci-lint
# make ci-lint
ci-lint: prepare-lint
	${GOPATH}/bin/golangci-lint run ./...

.PHONY: prepare-lint
# make prepare-lint
prepare-lint:
	@if ! which golangci-lint &>/dev/null; then \
  		echo "Installing golangci-lint"; \
  		curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s latest; \
  	fi

.PHONY: test
# make test
test:
	go test -cover ./... | grep -v vendor;true
	go vet ./... | grep -v vendor;true
	go test -short ${PKG_LIST}

.PHONY: cover
# make cover
cover:
	go test -short -coverprofile coverage.txt -covermode=atomic ${PKG_LIST}
	go tool cover -html=coverage.txt

.PHONY: docker
# make docker  生成docker镜像
docker:
	docker build -t eagle:$(versionDir) -f Dockeffile .

.PHONY: deploy
# make deploy  deploy app to k8s
deploy:
	sh deploy/deploy.sh

.PHONY: clean
# make clean
clean:
	@-rm -vrf eagle
	@-rm -vrf cover.out
	@-rm -vrf coverage.txt
	@go mod tidy
	@echo "clean finished"

.PHONY: docs
# gen swagger doc
docs:
	@if ! which swag &>/dev/null; then \
  		echo "downloading swag"; \
  		go get -u github.com/swaggo/swag/cmd/swag; \
  	fi
	@swag init
	@mv docs/docs.go api/http
	@mv docs/swagger.json api/http
	@mv docs/swagger.yaml api/http
	@echo "gen-docs done"
	@echo "see docs by: http://localhost:8080/swagger/index.html"

.PHONY: graph
# make graph 生成交互式的可视化Go程序调用图(会在浏览器自动打开)
graph:
	@export GO111MODULE="on"
	@if ! which go-callvis &>/dev/null; then \
  		echo "downloading go-callvis"; \
  		go get -u github.com/ofabry/go-callvis; \
  	fi
	@echo "generating graph"
	@go-callvis github.com/go-eagle/eagle

.PHONY: mockgen
# make mockgen gen mock file
mockgen:
	cd ./internal &&  for file in `egrep -rnl "type.*?interface" ./repository | grep -v "_test" `; do \
		echo $$file ; \
		cd .. && mockgen -destination="./internal/mock/$$file" -source="./internal/$$file" && cd ./internal ; \
	done

.PHONY: init
# init env
init:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go get -v github.com/google/gnostic
	go get -v github.com/google/gnostic/cmd/protoc-gen-openapi
	go get -u github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc
	go get github.com/golang/mock/mockgen

.PHONY: proto
# generate proto struct only
proto:
	protoc --proto_path=. \
           --proto_path=./third_party \
           --go_out=. --go_opt=paths=source_relative \
           $(API_PROTO_FILES)

.PHONY: grpc
# generate grpc code
grpc:
	protoc --proto_path=. \
           --proto_path=./third_party \
           --go_out=. --go_opt=paths=source_relative \
           --go-grpc_out=. --go-grpc_opt=paths=source_relative \
           $(API_PROTO_FILES)

.PHONY: openapi
# generate openapi
openapi:
	protoc --proto_path=. \
          --proto_path=./third_party \
          --openapi_out=. \
          $(API_PROTO_FILES)
	  
.PHONY: doc
# generate html or markdown doc
doc:
	protoc --proto_path=. \
           --proto_path=./third_party \
	   --doc_out=. \
	   --doc_opt=html,index.html \
	   $(API_PROTO_FILES)

# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m  %-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := all
