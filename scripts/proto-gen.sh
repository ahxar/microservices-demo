#!/bin/bash

set -e

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}Starting protobuf code generation...${NC}"

# Proto directory
PROTO_DIR="proto"

# Check if protoc is installed
if ! command -v protoc &> /dev/null; then
    echo "Error: protoc is not installed"
    echo "Install with: brew install protobuf (macOS) or apt-get install protobuf-compiler (Linux)"
    exit 1
fi

# Check if protoc-gen-go is installed
if ! command -v protoc-gen-go &> /dev/null; then
    echo "Error: protoc-gen-go is not installed"
    echo "Install with: go install google.golang.org/protobuf/cmd/protoc-gen-go@latest"
    exit 1
fi

# Check if protoc-gen-go-grpc is installed
if ! command -v protoc-gen-go-grpc &> /dev/null; then
    echo "Error: protoc-gen-go-grpc is not installed"
    echo "Install with: go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest"
    exit 1
fi

# Generate common proto
echo -e "${GREEN}Generating common proto...${NC}"
protoc \
    --proto_path=${PROTO_DIR} \
    --go_out=${PROTO_DIR} \
    --go_opt=paths=source_relative \
    ${PROTO_DIR}/common/v1/common.proto

# Generate user proto
echo -e "${GREEN}Generating user proto...${NC}"
protoc \
    --proto_path=${PROTO_DIR} \
    --go_out=${PROTO_DIR} \
    --go_opt=paths=source_relative \
    --go-grpc_out=${PROTO_DIR} \
    --go-grpc_opt=paths=source_relative \
    ${PROTO_DIR}/user/v1/user.proto

# Generate catalog proto
echo -e "${GREEN}Generating catalog proto...${NC}"
protoc \
    --proto_path=${PROTO_DIR} \
    --go_out=${PROTO_DIR} \
    --go_opt=paths=source_relative \
    --go-grpc_out=${PROTO_DIR} \
    --go-grpc_opt=paths=source_relative \
    ${PROTO_DIR}/catalog/v1/catalog.proto

# Generate cart proto
echo -e "${GREEN}Generating cart proto...${NC}"
protoc \
    --proto_path=${PROTO_DIR} \
    --go_out=${PROTO_DIR} \
    --go_opt=paths=source_relative \
    --go-grpc_out=${PROTO_DIR} \
    --go-grpc_opt=paths=source_relative \
    ${PROTO_DIR}/cart/v1/cart.proto

# Generate order proto
echo -e "${GREEN}Generating order proto...${NC}"
protoc \
    --proto_path=${PROTO_DIR} \
    --go_out=${PROTO_DIR} \
    --go_opt=paths=source_relative \
    --go-grpc_out=${PROTO_DIR} \
    --go-grpc_opt=paths=source_relative \
    ${PROTO_DIR}/order/v1/order.proto

# Generate payment proto
echo -e "${GREEN}Generating payment proto...${NC}"
protoc \
    --proto_path=${PROTO_DIR} \
    --go_out=${PROTO_DIR} \
    --go_opt=paths=source_relative \
    --go-grpc_out=${PROTO_DIR} \
    --go-grpc_opt=paths=source_relative \
    ${PROTO_DIR}/payment/v1/payment.proto

# Generate shipping proto
echo -e "${GREEN}Generating shipping proto...${NC}"
protoc \
    --proto_path=${PROTO_DIR} \
    --go_out=${PROTO_DIR} \
    --go_opt=paths=source_relative \
    --go-grpc_out=${PROTO_DIR} \
    --go-grpc_opt=paths=source_relative \
    ${PROTO_DIR}/shipping/v1/shipping.proto

# Generate notification proto
echo -e "${GREEN}Generating notification proto...${NC}"
protoc \
    --proto_path=${PROTO_DIR} \
    --go_out=${PROTO_DIR} \
    --go_opt=paths=source_relative \
    --go-grpc_out=${PROTO_DIR} \
    --go-grpc_opt=paths=source_relative \
    ${PROTO_DIR}/notification/v1/notification.proto

echo -e "${BLUE}Protobuf code generation completed!${NC}"
