#!/bin/bash

# Build script for ProxyPal NVIDIA Load Balancer

set -e

echo "Building ProxyPal NVIDIA Load Balancer..."

# Clean previous builds
rm -f proxypal

# Build for current platform
go build -o proxypal ./cmd/proxypal

echo "Build complete! Binary: ./proxypal"

# Make it executable
chmod +x proxypal

echo ""
echo "To run the application:"
echo "  1. Copy config.example.yaml to config.yaml"
echo "  2. Edit config.yaml and add your NVIDIA API keys"
echo "  3. Run: ./proxypal"
