#!/bin/bash

set -e

SWAGGER_UI_VERSION="4.18.3"
SWAGGER_UI_DIR="cmd/gateway/swagger-ui"

# Create directory if it doesn't exist
mkdir -p $SWAGGER_UI_DIR

# Download and extract Swagger UI
if [ ! -f "${SWAGGER_UI_DIR}/index.html" ]; then
  echo "Downloading Swagger UI v${SWAGGER_UI_VERSION}..."
  
  # Create temporary directory
  TMP_DIR=$(mktemp -d)
  
  # Download Swagger UI
  curl -sSL "https://github.com/swagger-api/swagger-ui/archive/v${SWAGGER_UI_VERSION}.tar.gz" -o "${TMP_DIR}/swagger-ui.tar.gz"
  
  # Extract files
  tar -xzf "${TMP_DIR}/swagger-ui.tar.gz" -C "${TMP_DIR}"
  
  # Copy only the dist files
  cp -r "${TMP_DIR}/swagger-ui-${SWAGGER_UI_VERSION}/dist"/* "${SWAGGER_UI_DIR}/"
  
  # Clean up temporary directory
  rm -rf "${TMP_DIR}"
  
  # Update the Swagger UI to use our local swagger.json
  sed -i '' 's#https://petstore.swagger.io/v2/swagger.json#/swagger.json#g' "${SWAGGER_UI_DIR}/swagger-initializer.js"
  
  echo "Swagger UI downloaded and configured successfully!"
else
  echo "Swagger UI already exists."
fi 