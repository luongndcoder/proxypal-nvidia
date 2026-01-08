# ProxyPal NVIDIA - Quick Deploy Guide

## 🚀 1-Minute Quick Start

```bash
# 1. Clone
git clone https://github.com/luongndcoder/proxypal-nvidia.git
cd proxypal-nvidia

# 2. Configure
cp config.example.yaml config.yaml
# Edit config.yaml and add your NVIDIA API keys

# 3. Deploy
./deploy.sh
```

That's it! Your load balancer is running on `http://localhost:8080`

## ⚡ Quick Commands

### Docker Deployment
```bash
./deploy.sh              # Auto build and deploy
docker-compose logs -f   # View logs
docker-compose down      # Stop service
```

### Local Development
```bash
make build              # Build binary
make run                # Run locally
make test               # Run tests
```

### Testing
```bash
./test.sh               # Run API tests
python example.py       # Python example
curl http://localhost:8080/health  # Health check
```

## 📊 Monitoring

```bash
# Statistics
curl http://localhost:8080/stats | jq

# Example response:
{
  "keys": 3,
  "stats": [
    {
      "KeyPrefix": "nvapi-...abc1",
      "RequestCount": 150,
      "ErrorCount": 0,
      "AvailableTokens": 38,
      "LastUsed": "2024-01-08T10:30:00Z"
    }
  ]
}
```

## 🔧 Configuration

Edit `config.yaml`:

```yaml
nvidia:
  api_keys:
    - "nvapi-YOUR-FIRST-KEY"
    - "nvapi-YOUR-SECOND-KEY"
    - "nvapi-YOUR-THIRD-KEY"
  rate_limit: 40  # Requests per minute per key
```

## 💻 Usage Example

### Python
```python
from openai import OpenAI

client = OpenAI(
    base_url="http://localhost:8080/v1",
    api_key="dummy"
)

response = client.chat.completions.create(
    model="minimaxai/minimax-m2",
    messages=[{"role": "user", "content": "Hello!"}],
    stream=True
)

for chunk in response:
    if chunk.choices[0].delta.content:
        print(chunk.choices[0].delta.content, end="")
```

### cURL
```bash
curl -X POST http://localhost:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "minimaxai/minimax-m2",
    "messages": [{"role": "user", "content": "Hello!"}],
    "stream": true
  }'
```

## 📚 Full Documentation

- [Complete README](README.md)
- [Quick Start Guide](QUICKSTART.md)
- [Project Structure](STRUCTURE.md)
- [Contributing](CONTRIBUTING.md)

## ❓ Common Issues

**Port 8080 already in use?**
```yaml
# Change port in config.yaml
server:
  port: 3000
```

**Rate limit errors?**
- Check `/stats` endpoint
- Add more API keys
- Wait for token refill

**Connection refused?**
```bash
# Check service status
docker-compose ps

# View logs
docker-compose logs
```

## 🎯 Features

- ✅ Round-robin load balancing
- ✅ 40 req/min per key rate limiting
- ✅ Auto-failover
- ✅ OpenAI compatible
- ✅ Streaming support
- ✅ Real-time statistics
- ✅ Docker ready
- ✅ Single binary

## 📦 What's Included

```
✓ Load balancer core (Go)
✓ Docker deployment
✓ Configuration files
✓ Test scripts
✓ Python examples
✓ Full documentation
```

## 🌟 Get NVIDIA API Keys

1. Visit https://build.nvidia.com/
2. Sign up for free
3. Get API key from any model page
4. Repeat for multiple accounts

## 🛠️ Build Info

- **Language**: Go 1.21
- **Binary Size**: ~12MB
- **Docker Image**: ~20MB
- **Dependencies**: 2 (gin, yaml)

---

**Need help?** Open an [issue](https://github.com/luongndcoder/proxypal-nvidia/issues)

**Like it?** Star the repo ⭐
