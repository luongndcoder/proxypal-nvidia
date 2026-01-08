# Ubuntu Server Deployment Guide

## ✅ Prerequisites

```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Install Docker Compose
sudo apt install docker-compose-plugin -y

# Add user to docker group (optional, to run without sudo)
sudo usermod -aG docker $USER
newgrp docker

# Verify installation
docker --version
docker compose version
```

## 🚀 Quick Deploy on Ubuntu

### Method 1: Using deploy script (Recommended)

```bash
# Clone repository
git clone https://github.com/luongndcoder/proxypal-nvidia.git
cd proxypal-nvidia

# Create config
cp config.example.yaml config.yaml

# Edit config and add your API keys
nano config.yaml
# or
vim config.yaml

# Deploy (auto-detects docker-compose or docker compose)
chmod +x deploy.sh
./deploy.sh
```

### Method 2: Manual Docker Compose

```bash
# Build image
docker build -t proxypal-nvidia:latest .

# Start service
docker compose up -d

# Check status
docker compose ps

# View logs
docker compose logs -f
```

### Method 3: Direct Docker Run

```bash
# Build
docker build -t proxypal-nvidia:latest .

# Run
docker run -d \
  --name proxypal \
  -p 8080:8080 \
  -v $(pwd)/config.yaml:/root/config.yaml \
  -e CONFIG_PATH=/root/config.yaml \
  --restart unless-stopped \
  proxypal-nvidia:latest
```

## 🔒 Security on Ubuntu

### 1. Firewall Configuration

```bash
# Allow port 8080
sudo ufw allow 8080/tcp

# Or if using nginx reverse proxy
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# Enable firewall
sudo ufw enable
sudo ufw status
```

### 2. Secure Config File

```bash
# Set correct permissions
chmod 600 config.yaml
chown $USER:$USER config.yaml
```

### 3. Run as non-root (Recommended)

Update `docker-compose.yml`:

```yaml
services:
  proxypal:
    user: "1000:1000"  # Your user UID:GID
    # ... rest of config
```

## 🌐 Nginx Reverse Proxy (Optional)

### Install Nginx

```bash
sudo apt install nginx -y
```

### Configure Nginx

Create `/etc/nginx/sites-available/proxypal`:

```nginx
server {
    listen 80;
    server_name your-domain.com;  # Change this

    location / {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;

        # For streaming support
        proxy_buffering off;
        proxy_read_timeout 300s;
    }
}
```

### Enable Site

```bash
sudo ln -s /etc/nginx/sites-available/proxypal /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx
```

### Add SSL (Let's Encrypt)

```bash
# Install certbot
sudo apt install certbot python3-certbot-nginx -y

# Get certificate
sudo certbot --nginx -d your-domain.com

# Auto-renewal is enabled by default
sudo certbot renew --dry-run
```

## 📊 Monitoring on Ubuntu

### 1. View Logs

```bash
# Docker logs
docker compose logs -f

# Follow specific service
docker compose logs -f proxypal

# Last 100 lines
docker compose logs --tail=100
```

### 2. Resource Monitoring

```bash
# Docker stats
docker stats proxypal

# System resources
htop  # Install: sudo apt install htop
```

### 3. Service Status

```bash
# Check if running
docker compose ps

# Check health
curl http://localhost:8080/health

# Check stats
curl http://localhost:8080/stats | jq
```

## 🔄 Auto-start on Reboot

### Method 1: Docker Compose (Recommended)

Already configured with `restart: unless-stopped` in docker-compose.yml

### Method 2: Systemd Service

Create `/etc/systemd/system/proxypal.service`:

```ini
[Unit]
Description=ProxyPal NVIDIA Load Balancer
Requires=docker.service
After=docker.service

[Service]
Type=oneshot
RemainAfterExit=yes
WorkingDirectory=/path/to/proxypal-nvidia
ExecStart=/usr/bin/docker compose up -d
ExecStop=/usr/bin/docker compose down
TimeoutStartSec=0

[Install]
WantedBy=multi-user.target
```

Enable service:

```bash
sudo systemctl daemon-reload
sudo systemctl enable proxypal
sudo systemctl start proxypal
sudo systemctl status proxypal
```

## 🛠️ Maintenance

### Update Container

```bash
# Pull latest code
git pull

# Rebuild and restart
docker compose down
docker compose build --no-cache
docker compose up -d
```

### Backup Configuration

```bash
# Backup config
cp config.yaml config.yaml.backup

# Backup with timestamp
cp config.yaml config.yaml.$(date +%Y%m%d_%H%M%S)
```

### Clean Up

```bash
# Remove stopped containers
docker container prune -f

# Remove unused images
docker image prune -f

# Remove all unused data
docker system prune -af
```

## 🐛 Troubleshooting Ubuntu

### Port Already in Use

```bash
# Check what's using port 8080
sudo lsof -i :8080
sudo netstat -tulpn | grep 8080

# Kill process
sudo kill -9 <PID>

# Or change port in config.yaml
```

### Permission Denied

```bash
# Add user to docker group
sudo usermod -aG docker $USER
newgrp docker

# Fix config file permissions
chmod 600 config.yaml
```

### Container Won't Start

```bash
# Check logs
docker compose logs

# Check Docker daemon
sudo systemctl status docker

# Restart Docker
sudo systemctl restart docker
```

### Out of Disk Space

```bash
# Check disk usage
df -h

# Docker disk usage
docker system df

# Clean up
docker system prune -af --volumes
```

## 📈 Performance Tuning

### Increase File Descriptors

Edit `/etc/security/limits.conf`:

```
* soft nofile 65536
* hard nofile 65536
```

### Docker Resource Limits

Update `docker-compose.yml`:

```yaml
services:
  proxypal:
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 512M
        reservations:
          cpus: '0.5'
          memory: 128M
```

## 🔍 Health Checks

### Create Health Check Script

`/usr/local/bin/proxypal-health.sh`:

```bash
#!/bin/bash
response=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/health)
if [ $response -eq 200 ]; then
    echo "OK"
    exit 0
else
    echo "FAIL"
    exit 1
fi
```

### Add to Cron

```bash
# Run health check every 5 minutes
crontab -e

# Add line:
*/5 * * * * /usr/local/bin/proxypal-health.sh || docker compose -f /path/to/proxypal-nvidia/docker-compose.yml restart
```

## 📱 Access from Outside

### Bind to All Interfaces

In `config.yaml`:

```yaml
server:
  host: "0.0.0.0"  # Allow external access
  port: 8080
```

### Configure Firewall

```bash
# Allow from specific IP
sudo ufw allow from 192.168.1.0/24 to any port 8080

# Or allow from anywhere (use with caution)
sudo ufw allow 8080/tcp
```

## ✅ Production Checklist

- [ ] Docker and Docker Compose installed
- [ ] Config file created with real API keys
- [ ] Config file permissions secured (chmod 600)
- [ ] Firewall configured
- [ ] Nginx reverse proxy setup (optional)
- [ ] SSL certificate installed (if using nginx)
- [ ] Auto-start on reboot enabled
- [ ] Monitoring/logging configured
- [ ] Backup strategy in place
- [ ] Health checks running

## 🚀 One-Line Deploy

```bash
curl -sSL https://raw.githubusercontent.com/luongndcoder/proxypal-nvidia/main/deploy.sh | bash
```

---

**Need help?** Open an issue on GitHub!
