# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial release of ProxyPal NVIDIA Load Balancer
- OpenAI-compatible API endpoints
- Round-robin load balancing across multiple API keys
- Token bucket rate limiting (40 requests/minute per key)
- Auto-failover when keys are rate limited
- Streaming and non-streaming response support
- Statistics endpoint for monitoring
- Health check endpoint
- Docker and docker-compose support
- Configuration via YAML file
- Comprehensive documentation and examples
- Unit tests for core functionality
- CI/CD with GitHub Actions

### Features

#### Load Balancing
- Round-robin distribution across API keys
- Automatic key rotation
- Request count tracking per key
- Error count tracking per key

#### Rate Limiting
- Token bucket algorithm
- Configurable rate limit per key
- Automatic token refill
- Real-time token availability tracking

#### API Endpoints
- `POST /v1/chat/completions` - OpenAI-compatible chat completions
- `GET /v1/models` - List available models
- `GET /health` - Health check
- `GET /stats` - Load balancer statistics

#### Configuration
- YAML-based configuration
- Environment variable support
- Multiple API key support
- Configurable rate limits
- Retry and failover settings
- Logging configuration

#### Documentation
- Comprehensive README
- Quick start guide
- API documentation
- Docker deployment guide
- Python and Node.js examples
- Contributing guidelines

## [1.0.0] - 2024-01-08

### Added
- Initial public release

---

## Version History

- **v1.0.0** - Initial release with core functionality
  - Load balancing
  - Rate limiting
  - OpenAI-compatible API
  - Docker support
  - Comprehensive documentation

---

## Planned Features

### v1.1.0
- [ ] Metrics export (Prometheus)
- [ ] Admin dashboard UI
- [ ] Multiple model endpoint support
- [ ] Request queue management
- [ ] Advanced retry strategies

### v1.2.0
- [ ] Database persistence for statistics
- [ ] User authentication
- [ ] API key management UI
- [ ] Webhook notifications
- [ ] Advanced logging (structured logging)

### v2.0.0
- [ ] Multi-provider support (OpenAI, Anthropic, etc.)
- [ ] Intelligent routing based on model capabilities
- [ ] Cost optimization features
- [ ] Request caching
- [ ] Load prediction and auto-scaling

---

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for details on how to contribute to this changelog.

## Links

- [Homepage](https://github.com/luongndcoder/proxypal-nvidia)
- [Issues](https://github.com/luongndcoder/proxypal-nvidia/issues)
- [Releases](https://github.com/luongndcoder/proxypal-nvidia/releases)
