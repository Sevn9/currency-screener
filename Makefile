PROTO_SRC=proto/currency/currency_service.proto
PROTO_OUT=.

PROTOC_GEN_GO=protoc-gen-go
PROTOC_GEN_GO_GRPC=protoc-gen-go-grpc

GEN_TEST_DATA_SCRIPT=./scripts/generate_test_data.go

K6_SCRIPT=k6-script.js

.PHONY: all build run test proto

all: build

install-tools:
	@echo "Checking and installing necessary tools..."
	@if ! [ -x "$$(command -v $(PROTOC_GEN_GO))" ]; then \
		echo "Installing protoc-gen-go..."; \
		go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.10; \
	fi
	@if ! [ -x "$$(command -v $(PROTOC_GEN_GO_GRPC))" ]; then \
		echo "Installing protoc-gen-go-grpc..."; \
		go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@1.5.1; \
	fi

proto: install-tools
	@echo "Generating gRPC and Protobuf code..."
	protoc --proto_path=proto --go_out=$(PROTO_OUT) --go-grpc_out=$(PROTO_OUT) $(PROTO_SRC)

build: proto
	go build -o bin/app

run: build
	./bin/app

test:
	go test -v ./... -cover

generate-test-data:
	@echo "Generating and inserting test data into the database..."
	go run $(GEN_TEST_DATA_SCRIPT)

mocks:
	go generate ./...

load-test:
	@echo "docker-compose up..."
	docker-compose up -d
	sleep 5
	@echo "Running load tests with k6..."
	k6 run $(K6_SCRIPT)