#!/bin/bash

# ProxyPal NVIDIA - Ubuntu Quick Deploy Script
# This script will install Docker, build and deploy ProxyPal on Ubuntu

set -e

echo "============================================================"
echo "  ProxyPal NVIDIA - Ubuntu Auto Deploy"
echo "============================================================"
echo ""

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if running on Linux
if [[ "$OSTYPE" != "linux-gnu"* ]]; then
    echo -e "${RED}❌ This script is designed for Ubuntu/Debian Linux${NC}"
    exit 1
fi

# Check if running as root
if [[ $EUID -eq 0 ]]; then
    echo -e "${YELLOW}⚠️  Running as root. It's recommended to run as regular user.${NC}"
    read -p "Continue anyway? (y/n) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Install Docker if not present
if ! command_exists docker; then
    echo -e "${YELLOW}🔧 Docker not found. Installing Docker...${NC}"
    curl -fsSL https://get.docker.com -o get-docker.sh
    sudo sh get-docker.sh
    rm get-docker.sh

    # Add current user to docker group
    sudo usermod -aG docker $USER
    echo -e "${GREEN}✅ Docker installed successfully!${NC}"
    echo -e "${YELLOW}⚠️  You may need to logout and login again for docker group to take effect${NC}"
else
    echo -e "${GREEN}✅ Docker is already installed${NC}"
fi

# Detect docker-compose command
if command_exists docker-compose; then
    COMPOSE_CMD="docker-compose"
    echo -e "${GREEN}✅ docker-compose is installed${NC}"
elif docker compose version &> /dev/null 2>&1; then
    COMPOSE_CMD="docker compose"
    echo -e "${GREEN}✅ Docker Compose (plugin) is installed${NC}"
else
    echo -e "${YELLOW}🔧 Installing docker-compose...${NC}"
    sudo apt-get update
    sudo apt-get install -y docker-compose
    COMPOSE_CMD="docker-compose"
    echo -e "${GREEN}✅ docker-compose installed successfully!${NC}"
fi

# Check if config.yaml exists
if [ ! -f "config.yaml" ]; then
    echo ""
    echo -e "${YELLOW}⚠️  config.yaml not found!${NC}"

    if [ -f "config.example.yaml" ]; then
        echo "Creating config.yaml from example..."
        cp config.example.yaml config.yaml
        chmod 600 config.yaml

        echo ""
        echo -e "${RED}⚠️  IMPORTANT: Edit config.yaml and add your NVIDIA API keys!${NC}"
        echo ""
        echo "Options to edit:"
        echo "  1. nano config.yaml"
        echo "  2. vim config.yaml"
        echo "  3. vi config.yaml"
        echo ""
        read -p "Would you like to edit now? (y/n) " -n 1 -r
        echo

        if [[ $REPLY =~ ^[Yy]$ ]]; then
            if command_exists nano; then
                nano config.yaml
            elif command_exists vim; then
                vim config.yaml
            elif command_exists vi; then
                vi config.yaml
            else
                echo -e "${RED}No text editor found. Please edit config.yaml manually.${NC}"
                exit 1
            fi
        else
            echo -e "${YELLOW}Please edit config.yaml before running the service!${NC}"
            exit 1
        fi
    else
        echo -e "${RED}config.example.yaml not found!${NC}"
        exit 1
    fi
else
    echo -e "${GREEN}✅ config.yaml found${NC}"
fi

# Verify config has real API keys
if grep -q "nvapi-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" config.yaml; then
    echo ""
    echo -e "${RED}⚠️  WARNING: config.yaml still contains example API keys!${NC}"
    echo "Please add your real NVIDIA API keys before deploying."
    read -p "Continue anyway? (y/n) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

echo ""
echo "🔨 Building Docker image..."
docker build -t proxypal-nvidia:latest . || {
    echo -e "${RED}❌ Failed to build Docker image${NC}"
    exit 1
}

echo ""
echo -e "${GREEN}✅ Docker image built successfully!${NC}"
echo ""

# Stop existing container if running
if docker ps -a | grep -q proxypal; then
    echo "🛑 Stopping existing ProxyPal container..."
    $COMPOSE_CMD down 2>/dev/null || docker stop proxypal 2>/dev/null || true
    docker rm proxypal 2>/dev/null || true
fi

echo "🚀 Starting ProxyPal..."
$COMPOSE_CMD up -d || {
    echo -e "${RED}❌ Failed to start with $COMPOSE_CMD${NC}"
    exit 1
}

# Wait for service to start
echo "⏳ Waiting for service to start..."
sleep 3

# Check if service is running
if $COMPOSE_CMD ps | grep -q "Up"; then
    echo ""
    echo -e "${GREEN}✅ ProxyPal is now running!${NC}"
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

    # Try to get server IP
    SERVER_IP=$(hostname -I | awk '{print $1}')
    if [ ! -z "$SERVER_IP" ]; then
        echo "  External URL: http://${SERVER_IP}:8080"
        echo "  (Make sure port 8080 is open in firewall)"
        echo ""
    fi

    # Test health endpoint
    if command_exists curl; then
        echo "🧪 Testing health endpoint..."
        if curl -s http://localhost:8080/health > /dev/null; then
            echo -e "${GREEN}✅ Health check passed!${NC}"
        else
            echo -e "${YELLOW}⚠️  Health check failed. Service may still be starting...${NC}"
        fi
        echo ""
    fi

    echo "📝 Useful commands:"
    echo "  View logs:      $COMPOSE_CMD logs -f"
    echo "  Stop service:   $COMPOSE_CMD down"
    echo "  Restart:        $COMPOSE_CMD restart"
    echo "  Check status:   $COMPOSE_CMD ps"
    echo ""

    # Configure firewall if ufw is available
    if command_exists ufw; then
        echo "🔒 Firewall Configuration:"
        if sudo ufw status | grep -q "Status: active"; then
            echo "UFW is active. You may need to allow port 8080:"
            echo "  sudo ufw allow 8080/tcp"
        fi
        echo ""
    fi

    echo "============================================================"
    echo ""
    echo -e "${GREEN}🎉 Deployment successful!${NC}"
    echo ""
    echo "Next steps:"
    echo "  1. Test: curl http://localhost:8080/health"
    echo "  2. View stats: curl http://localhost:8080/stats"
    echo "  3. Use with OpenAI client library"
    echo ""
    echo "Documentation: See README.md and UBUNTU_DEPLOY.md"
    echo ""
else
    echo -e "${RED}❌ Failed to start service${NC}"
    echo "Check logs with: $COMPOSE_CMD logs"
    exit 1
fi
