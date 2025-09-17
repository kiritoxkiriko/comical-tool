#!/bin/bash

# Script to generate IDL files for Hertz
# This script generates Go code from Thrift and Protocol Buffers definitions

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}🚀 Starting IDL generation for Comical Tool...${NC}"

# Check if required tools are installed
check_tool() {
    if ! command -v $1 &> /dev/null; then
        echo -e "${RED}❌ $1 is not installed. Please install it first.${NC}"
        exit 1
    fi
}

echo -e "${YELLOW}📋 Checking required tools...${NC}"
check_tool "thrift"
check_tool "protoc"

# Create output directories
mkdir -p api/generated/thrift
mkdir -p api/generated/proto

echo -e "${YELLOW}📝 Generating Thrift files...${NC}"
# Generate Thrift files
thrift --gen go:package_prefix=github.com/kiritoxkiriko/comical-tool/api/generated/thrift/ \
       --out api/generated/thrift/ \
       api/short_url.thrift

echo -e "${YELLOW}📝 Generating Protocol Buffer files...${NC}"
# Generate Protocol Buffer files
protoc --go_out=api/generated/proto/ \
       --go_opt=paths=source_relative \
       --go-grpc_out=api/generated/proto/ \
       --go-grpc_opt=paths=source_relative \
       api/short_url.proto

echo -e "${GREEN}✅ IDL generation completed successfully!${NC}"
echo -e "${GREEN}📁 Generated files:${NC}"
echo -e "  - Thrift: api/generated/thrift/"
echo -e "  - Protocol Buffers: api/generated/proto/"

# Update go.mod if needed
echo -e "${YELLOW}🔄 Updating dependencies...${NC}"
go mod tidy

echo -e "${GREEN}🎉 All done! You can now use the generated code in your Hertz application.${NC}"
