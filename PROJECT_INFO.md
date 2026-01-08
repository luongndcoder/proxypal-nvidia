# 🚀 ProxyPal NVIDIA Load Balancer - Project Information

## 📦 Project Overview

**ProxyPal** is a production-ready load balancer for NVIDIA API keys, designed to maximize free tier usage through intelligent request distribution and rate limiting.

### 🎯 Key Objectives

1. **Maximize Free Tier**: Distribute requests across multiple NVIDIA accounts
2. **Rate Limit Compliance**: Respect 40 requests/minute per key limit
3. **High Availability**: Auto-failover when keys are exhausted
4. **Developer Friendly**: OpenAI-compatible, drop-in replacement
5. **Production Ready**: Docker deployment, monitoring, logging

## 📊 Project Statistics

| Metric | Value |
|--------|-------|
| **Lines of Code** | 1,005 LOC (Go) |
| **Binary Size** | ~12 MB |
| **Docker Image** | ~20 MB |
| **Dependencies** | 2 (minimal) |
| **Test Coverage** | Core functionality |
| **Build Time** | ~10 seconds |
| **Startup Time** | <1 second |

## 🏗️ Architecture

```
┌─────────────┐
│   Client    │
│  (OpenAI)   │
└──────┬──────┘
       │
       ▼
┌─────────────────────────────────────┐
│     Gin Router (CORS, Routing)      │
└──────┬──────────────────────────────┘
       │
       ▼
┌─────────────────────────────────────┐
│    Proxy Server (Request Handler)   │
└──────┬──────────────────────────────┘
       │
       ▼
┌─────────────────────────────────────┐
│   Load Balancer (Round-Robin)       │
└──────┬──────────────────────────────┘
       │
       ▼
┌─────────────────────────────────────┐
│   Rate Limiter (Token Bucket)       │
└──────┬──────────────────────────────┘
       │
       ├──────┬──────┬──────┐
       ▼      ▼      ▼      ▼
    Key 1  Key 2  Key 3  Key N
       │      │      │      │
       └──────┴──────┴──────┘
              │
              ▼
       ┌─────────────┐
       │ NVIDIA API  │
       └─────────────┘
```

## 🛠️ Technology Stack

### Core
- **Language**: Go 1.21
- **Web Framework**: Gin (high-performance HTTP)
- **Config**: YAML (gopkg.in/yaml.v3)

### Development
- **Testing**: Go test framework
- **CI/CD**: GitHub Actions
- **Linting**: go vet, go fmt
- **Build**: Make, Docker

### Deployment
- **Container**: Docker + Docker Compose
- **Platform**: Any (Linux, macOS, Windows)
- **Cloud Ready**: AWS, GCP, Azure, etc.

## 📁 Project Structure

```
proxypal-nvidia/
├── cmd/proxypal/           # Application entry point
│   └── main.go
├── internal/               # Internal packages
│   ├── balancer/           # Load balancing logic
│   │   ├── loadbalancer.go
│   │   ├── loadbalancer_test.go
│   │   ├── ratelimiter.go
│   │   └── ratelimiter_test.go
│   ├── config/             # Configuration
│   │   ├── config.go
│   │   └── config_test.go
│   └── proxy/              # HTTP proxy
│       └── proxy.go
├── .github/workflows/      # CI/CD
│   └── ci.yml
├── docs/                   # Documentation
│   ├── README.md
│   ├── QUICKSTART.md
│   ├── DEPLOY.md
│   ├── STRUCTURE.md
│   ├── CONTRIBUTING.md
│   ├── CHANGELOG.md
│   └── SUMMARY.md
├── config.example.yaml     # Config template
├── docker-compose.yml      # Docker Compose
├── Dockerfile              # Docker image
├── Makefile                # Build automation
├── build.sh                # Build script
├── deploy.sh               # Deploy script
├── test.sh                 # Test script
└── example.py              # Usage example
```

## ✨ Features

### Core Features
- ✅ **Round-Robin Load Balancing**: Even distribution across keys
- ✅ **Rate Limiting**: Token bucket (40 req/min per key)
- ✅ **Auto-Failover**: Automatic key switching
- ✅ **Request Tracking**: Per-key statistics
- ✅ **Error Handling**: Graceful error management

### API Features
- ✅ **OpenAI Compatible**: Drop-in replacement
- ✅ **Streaming Support**: Server-sent events
- ✅ **Model Listing**: GET /v1/models
- ✅ **Health Check**: GET /health
- ✅ **Statistics**: GET /stats

### DevOps Features
- ✅ **Docker Support**: Containerized deployment
- ✅ **Single Binary**: No dependencies
- ✅ **YAML Config**: Easy configuration
- ✅ **Environment Variables**: 12-factor app
- ✅ **Logging**: Configurable logging

## 🔧 Configuration

### Minimal Config
```yaml
nvidia:
  api_keys:
    - "nvapi-key-1"
    - "nvapi-key-2"
```

### Full Config
```yaml
server:
  port: 8080
  host: "0.0.0.0"

nvidia:
  base_url: "https://integrate.api.nvidia.com/v1"
  rate_limit: 40
  api_keys: [...]
  timeout: 300
  retry:
    max_retries: 3
    auto_failover: true

logging:
  level: "info"
  enable_request_log: true
```

## 🚀 Deployment Options

### 1. Docker (Recommended)
```bash
./deploy.sh
```

### 2. Docker Compose
```bash
docker-compose up -d
```

### 3. Binary
```bash
make build && ./proxypal
```

### 4. Source
```bash
go run ./cmd/proxypal
```

## 📈 Performance Benchmarks

### With 3 API Keys
- **Throughput**: ~120 req/min (2 req/sec)
- **Latency**: ~500ms (network dependent)
- **Memory**: ~20-50 MB
- **CPU**: <5% (idle), <20% (peak)

### Scaling
- **5 keys**: ~200 req/min
- **10 keys**: ~400 req/min
- **20 keys**: ~800 req/min

## 🧪 Testing

### Unit Tests
```bash
go test ./... -v
```

### Integration Tests
```bash
./test.sh
```

### Load Testing
```bash
# Example with Apache Bench
ab -n 100 -c 10 -p request.json \
   -T application/json \
   http://localhost:8080/v1/chat/completions
```

## 📊 Monitoring

### Statistics Endpoint
```bash
curl http://localhost:8080/stats
```

### Response Example
```json
{
  "keys": 3,
  "stats": [
    {
      "KeyPrefix": "nvapi-...abc1",
      "RequestCount": 150,
      "ErrorCount": 2,
      "AvailableTokens": 38,
      "LastUsed": "2024-01-08T10:30:00Z"
    }
  ],
  "timestamp": "2024-01-08T10:30:05Z"
}
```

## 🔒 Security

### Implemented
- ✅ API key masking in logs
- ✅ No key storage in responses
- ✅ CORS support
- ✅ Input validation
- ✅ Secure config file handling

### Recommendations
- 🔐 Use HTTPS in production
- 🔐 Secure config.yaml permissions (chmod 600)
- 🔐 Rotate API keys regularly
- 🔐 Use reverse proxy (nginx, Caddy)
- 🔐 Monitor usage patterns

## 🎯 Use Cases

1. **Development & Testing**
   - Free tier LLM access
   - Prototype development
   - CI/CD integration

2. **Small Projects**
   - Personal projects
   - MVPs
   - Educational use

3. **Load Distribution**
   - High availability
   - Cost optimization
   - Rate limit avoidance

## 🗺️ Roadmap

### v1.1.0 (Planned)
- [ ] Prometheus metrics
- [ ] Web dashboard
- [ ] Request queue
- [ ] Enhanced retry logic

### v1.2.0 (Planned)
- [ ] Database persistence
- [ ] User authentication
- [ ] Admin API
- [ ] Webhook notifications

### v2.0.0 (Future)
- [ ] Multi-provider support
- [ ] Intelligent routing
- [ ] Cost optimization
- [ ] Request caching

## 📚 Documentation

| Document | Description |
|----------|-------------|
| [README.md](README.md) | Complete documentation |
| [DEPLOY.md](DEPLOY.md) | Quick deployment guide |
| [QUICKSTART.md](QUICKSTART.md) | Getting started |
| [STRUCTURE.md](STRUCTURE.md) | Project structure |
| [CONTRIBUTING.md](CONTRIBUTING.md) | How to contribute |
| [CHANGELOG.md](CHANGELOG.md) | Version history |
| [SUMMARY.md](SUMMARY.md) | Project summary |

## 🤝 Contributing

We welcome contributions! See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

### Quick Contribution Guide
1. Fork the repo
2. Create feature branch
3. Make changes
4. Add tests
5. Submit PR

## 📄 License

MIT License - Free for personal and commercial use.

## 🙏 Acknowledgments

- Inspired by [CLIProxyAPI](https://github.com/router-for-me/CLIProxyAPI)
- Built with [Gin Web Framework](https://gin-gonic.com/)
- NVIDIA for free tier API access
- Community contributors

## 📞 Support

- **Issues**: [GitHub Issues](https://github.com/luongndcoder/proxypal-nvidia/issues)
- **Discussions**: [GitHub Discussions](https://github.com/luongndcoder/proxypal-nvidia/discussions)
- **Documentation**: This repository

## 🌟 Star History

If you find this project useful, please star it! ⭐

---

**Built with ❤️ for the developer community**

Last Updated: 2024-01-08
