#!/bin/bash

# Test script for ProxyPal NVIDIA Load Balancer

BASE_URL="http://localhost:8080"

echo "=================================================="
echo "ProxyPal NVIDIA Load Balancer - Test Suite"
echo "=================================================="
echo ""

# Test 1: Health check
echo "Test 1: Health Check"
echo "--------------------"
curl -s "${BASE_URL}/health" | jq '.' || echo "Error: Failed to connect"
echo ""
echo ""

# Test 2: Statistics
echo "Test 2: Statistics"
echo "------------------"
curl -s "${BASE_URL}/stats" | jq '.' || echo "Error: Failed to get stats"
echo ""
echo ""

# Test 3: List models
echo "Test 3: List Models"
echo "-------------------"
curl -s "${BASE_URL}/v1/models" | jq '.data[0:3]' || echo "Error: Failed to list models"
echo ""
echo ""

# Test 4: Non-streaming chat completion
echo "Test 4: Non-Streaming Chat Completion"
echo "--------------------------------------"
curl -s -X POST "${BASE_URL}/v1/chat/completions" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "minimaxai/minimax-m2",
    "messages": [{"role": "user", "content": "Say hello in one sentence"}],
    "max_tokens": 50,
    "stream": false
  }' | jq '.choices[0].message.content' || echo "Error: Failed to complete chat"
echo ""
echo ""

# Test 5: Streaming chat completion
echo "Test 5: Streaming Chat Completion"
echo "----------------------------------"
echo "Response: "
curl -s -X POST "${BASE_URL}/v1/chat/completions" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "minimaxai/minimax-m2",
    "messages": [{"role": "user", "content": "Count from 1 to 5"}],
    "max_tokens": 50,
    "stream": true
  }' || echo "Error: Failed to stream"
echo ""
echo ""

echo "=================================================="
echo "All tests completed!"
echo "=================================================="
