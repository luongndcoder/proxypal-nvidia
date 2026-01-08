#!/bin/bash

# Verification script for ProxyPal deployment

echo "============================================================"
echo "  ProxyPal NVIDIA - Deployment Verification"
echo "============================================================"
echo ""

# Detect compose command
if command -v docker-compose &> /dev/null; then
    COMPOSE_CMD="docker-compose"
elif docker compose version &> /dev/null 2>&1; then
    COMPOSE_CMD="docker compose"
else
    echo "❌ docker-compose not found"
    exit 1
fi

echo "Using: $COMPOSE_CMD"
echo ""

# Check if container is running
echo "1. Checking container status..."
if $COMPOSE_CMD ps | grep -q "Up"; then
    echo "✅ Container is running"
else
    echo "❌ Container is not running"
    echo "Start with: $COMPOSE_CMD up -d"
    exit 1
fi
echo ""

# Check health endpoint
echo "2. Checking health endpoint..."
if curl -s http://localhost:8080/health > /dev/null 2>&1; then
    response=$(curl -s http://localhost:8080/health)
    echo "✅ Health check passed"
    echo "   Response: $response"
else
    echo "❌ Health check failed"
    echo "   Make sure the service is running on port 8080"
fi
echo ""

# Check stats endpoint
echo "3. Checking stats endpoint..."
if curl -s http://localhost:8080/stats > /dev/null 2>&1; then
    echo "✅ Stats endpoint is accessible"
    if command -v jq &> /dev/null; then
        echo "   Stats:"
        curl -s http://localhost:8080/stats | jq '.stats[] | {KeyPrefix, RequestCount, AvailableTokens}'
    fi
else
    echo "❌ Stats endpoint failed"
fi
echo ""

# Check models endpoint
echo "4. Checking models endpoint..."
if curl -s http://localhost:8080/v1/models > /dev/null 2>&1; then
    echo "✅ Models endpoint is accessible"
else
    echo "❌ Models endpoint failed"
fi
echo ""

# Check logs for errors
echo "5. Checking recent logs..."
if $COMPOSE_CMD logs --tail=20 2>&1 | grep -i error; then
    echo "⚠️  Found errors in logs (see above)"
else
    echo "✅ No errors in recent logs"
fi
echo ""

# Summary
echo "============================================================"
echo "  Summary"
echo "============================================================"
echo "  Service: http://localhost:8080"
echo "  Health: http://localhost:8080/health"
echo "  Stats: http://localhost:8080/stats"
echo ""
echo "  Logs: $COMPOSE_CMD logs -f"
echo "  Stop: $COMPOSE_CMD down"
echo "============================================================"
