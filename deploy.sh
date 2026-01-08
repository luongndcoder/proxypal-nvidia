#!/bin/bash

# Quick deployment script for ProxyPal NVIDIA Load Balancer

set -e

echo "============================================================"
echo "  ProxyPal NVIDIA Load Balancer - Quick Deployment"
echo "============================================================"
echo ""

# Check if config.yaml exists
if [ ! -f "config.yaml" ]; then
    echo "⚠️  config.yaml not found!"
    echo ""
    echo "Creating config.yaml from example..."
    cp config.example.yaml config.yaml
    echo ""
    echo "✅ config.yaml created!"
    echo ""
    echo "⚠️  IMPORTANT: Please edit config.yaml and add your NVIDIA API keys:"
    echo "   1. Open config.yaml in your editor"
    echo "   2. Find the 'api_keys' section under 'nvidia'"
    echo "   3. Replace the example keys with your real NVIDIA API keys"
    echo ""
    echo "Then run this script again to deploy."
    exit 1
fi

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed!"
    echo "Please install Docker first: https://docs.docker.com/get-docker/"
    exit 1
fi

# Check if docker-compose is available
if command -v docker-compose &> /dev/null; then
    COMPOSE_CMD="docker-compose"
elif docker compose version &> /dev/null 2>&1; then
    COMPOSE_CMD="docker compose"
else
    echo "❌ docker-compose is not available!"
    echo "Please install docker-compose: https://docs.docker.com/compose/install/"
    exit 1
fi

echo "🔨 Building Docker image..."
docker build -t proxypal-nvidia:latest . || {
    echo "❌ Failed to build Docker image"
    exit 1
}

echo ""
echo "✅ Docker image built successfully!"
echo ""

echo "🚀 Starting ProxyPal with $COMPOSE_CMD..."
$COMPOSE_CMD up -d || {
    echo "❌ Failed to start with $COMPOSE_CMD"
    exit 1
}

echo ""
echo "✅ ProxyPal is now running!"
echo ""
echo "============================================================"
echo "  Service Information"
echo "============================================================"
echo "  URL:        http://localhost:8080"
echo "  Health:     http://localhost:8080/health"
echo "  Stats:      http://localhost:8080/stats"
echo "  Models:     http://localhost:8080/v1/models"
echo "============================================================"
echo ""
echo "📝 Useful commands:"
echo "  View logs:      $COMPOSE_CMD logs -f"
echo "  Stop service:   $COMPOSE_CMD down"
echo "  Restart:        $COMPOSE_CMD restart"
echo "  Check status:   $COMPOSE_CMD ps"
echo ""
echo "🧪 Test the service:"
echo "  curl http://localhost:8080/health"
echo "  curl http://localhost:8080/stats"
echo ""
echo "Happy load balancing! 🚀"
