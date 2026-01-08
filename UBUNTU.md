# 🚀 ProxyPal NVIDIA - Ubuntu Quick Start

## ⚡ 1-Command Deploy

```bash
curl -sSL https://raw.githubusercontent.com/luongndcoder/proxypal-nvidia/main/ubuntu-deploy.sh | bash
```

## 📦 Manual Install

### Step 1: Clone
```bash
git clone https://github.com/luongndcoder/proxypal-nvidia.git
cd proxypal-nvidia
```

### Step 2: Configure
```bash
cp config.example.yaml config.yaml
nano config.yaml  # Add your NVIDIA API keys
```

### Step 3: Deploy
```bash
chmod +x ubuntu-deploy.sh
./ubuntu-deploy.sh
```

## 🐳 What Gets Installed

- ✅ Docker (if not installed)
- ✅ docker-compose or Docker Compose plugin
- ✅ ProxyPal container (auto-start on reboot)

## 🔧 Commands

```bash
# Using docker-compose (old)
docker-compose logs -f
docker-compose down
docker-compose restart

# Using docker compose (new)
docker compose logs -f
docker compose down
docker compose restart

# Scripts auto-detect which one to use
./deploy.sh      # Works with both
./verify.sh      # Check if everything is working
```

## ✅ Verify Deployment

```bash
./verify.sh

# Or manually
curl http://localhost:8080/health
curl http://localhost:8080/stats
```

## 🌐 Access from Outside

### Option 1: Direct Access
```bash
# Open firewall
sudo ufw allow 8080/tcp

# Access via IP
curl http://YOUR_SERVER_IP:8080/health
```

### Option 2: Nginx Reverse Proxy (Recommended)
See [UBUNTU_DEPLOY.md](UBUNTU_DEPLOY.md) for full nginx setup with SSL.

## 📊 Monitor

```bash
# Logs
docker-compose logs -f  # or: docker compose logs -f

# Stats
curl http://localhost:8080/stats | jq

# System resources
docker stats proxypal-nvidia-proxypal-1
```

## 🔄 Update

```bash
git pull
docker-compose down
docker-compose build --no-cache
docker-compose up -d
```

## 🛠️ Troubleshooting

### Container won't start
```bash
# Check logs
docker-compose logs

# Check if port 8080 is free
sudo lsof -i :8080

# Restart Docker
sudo systemctl restart docker
```

### Permission denied
```bash
# Add user to docker group
sudo usermod -aG docker $USER
newgrp docker
```

### Config file issues
```bash
# Check permissions
ls -la config.yaml

# Fix permissions
chmod 600 config.yaml
```

## 📚 Full Documentation

- [Complete Guide](README.md)
- [Ubuntu Deploy Guide](UBUNTU_DEPLOY.md)
- [Quick Start](QUICKSTART.md)

## 🎯 Quick Test

```bash
# Install Python OpenAI library
pip install openai

# Run example
python example.py
```

---

**Support**: [GitHub Issues](https://github.com/luongndcoder/proxypal-nvidia/issues)
