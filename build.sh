#!/bin/bash

# Build script for ProxyPal NVIDIA Load Balancer

set -e

echo "============================================================"
echo "  Building ProxyPal NVIDIA Load Balancer"
echo "============================================================"
echo ""

# Clean previous builds
echo "🧹 Cleaning previous builds..."
rm -f proxypal
echo ""

# Build for current platform
echo "🔨 Building for current platform..."
go build -o proxypal ./cmd/proxypal

# Make it executable
chmod +x proxypal

# Get file size
SIZE=$(ls -lh proxypal | awk '{print $5}')

echo ""
echo "============================================================"
echo "  Build Complete!"
echo "============================================================"
echo "  Binary: ./proxypal"
echo "  Size: $SIZE"
echo "============================================================"
echo ""
echo "To run the application:"
echo "  1. Copy config.example.yaml to config.yaml"
echo "  2. Edit config.yaml and add your NVIDIA API keys"
echo "  3. Run: ./proxypal"
echo ""
echo "Or use: make run"
