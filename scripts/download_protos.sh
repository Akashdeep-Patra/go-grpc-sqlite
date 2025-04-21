#!/bin/bash

set -e

PROTO_DIR="./third_party/proto"
mkdir -p $PROTO_DIR

# Download Google APIs proto files
if [ ! -d "${PROTO_DIR}/google" ]; then
  echo "Downloading Google APIs proto files..."
  mkdir -p ${PROTO_DIR}/google/api
  
  curl -sSL "https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/annotations.proto" \
    -o "${PROTO_DIR}/google/api/annotations.proto"
  curl -sSL "https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/http.proto" \
    -o "${PROTO_DIR}/google/api/http.proto"
  curl -sSL "https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/field_behavior.proto" \
    -o "${PROTO_DIR}/google/api/field_behavior.proto"
fi

# Download Protoc Gen OpenAPI V2 proto files
if [ ! -d "${PROTO_DIR}/protoc-gen-openapiv2" ]; then
  echo "Downloading Protoc Gen OpenAPI V2 proto files..."
  mkdir -p ${PROTO_DIR}/protoc-gen-openapiv2/options
  
  curl -sSL "https://raw.githubusercontent.com/grpc-ecosystem/grpc-gateway/main/protoc-gen-openapiv2/options/annotations.proto" \
    -o "${PROTO_DIR}/protoc-gen-openapiv2/options/annotations.proto"
  curl -sSL "https://raw.githubusercontent.com/grpc-ecosystem/grpc-gateway/main/protoc-gen-openapiv2/options/openapiv2.proto" \
    -o "${PROTO_DIR}/protoc-gen-openapiv2/options/openapiv2.proto"
fi

echo "All required proto files downloaded successfully!" 