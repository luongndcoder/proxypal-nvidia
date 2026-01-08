# Project Structure

```
proxypal-nvidia/
├── cmd/
│   └── proxypal/
│       └── main.go              # Application entry point
│
├── internal/
│   ├── balancer/
│   │   ├── loadbalancer.go      # Round-robin load balancing logic
│   │   └── ratelimiter.go       # Token bucket rate limiter
│   ├── config/
│   │   └── config.go            # Configuration management
│   └── proxy/
│       └── proxy.go             # HTTP proxy server & handlers
│
├── config.example.yaml          # Example configuration file
├── docker-compose.yml           # Docker Compose configuration
├── Dockerfile                   # Docker image definition
├── Makefile                     # Build automation
├── build.sh                     # Build script
├── test.sh                      # Test script
├── example.py                   # Python usage example
├── .gitignore                   # Git ignore rules
├── .env.example                 # Environment variables example
├── LICENSE                      # MIT License
├── README.md                    # Main documentation
├── QUICKSTART.md               # Quick start guide
└── go.mod                       # Go module definition
```

## Core Components

### cmd/proxypal/main.go
- Application entry point
- Server initialization
- Route setup
- CORS middleware
- Startup banner

### internal/balancer/
- **loadbalancer.go**: Manages multiple API keys with round-robin distribution
- **ratelimiter.go**: Token bucket algorithm for rate limiting (40 req/min per key)

### internal/config/
- **config.go**: Configuration loading and validation from YAML

### internal/proxy/
- **proxy.go**: HTTP proxy server with OpenAI-compatible endpoints
  - POST /v1/chat/completions (streaming & non-streaming)
  - GET /v1/models
  - GET /health
  - GET /stats

## Configuration Files

- **config.example.yaml**: Template configuration with all available options
- **docker-compose.yml**: Docker deployment configuration
- **Dockerfile**: Multi-stage Docker build
- **.env.example**: Environment variables template

## Scripts

- **build.sh**: Simple build script
- **test.sh**: API testing script
- **example.py**: Python client example
- **Makefile**: Advanced build automation

## Documentation

- **README.md**: Complete project documentation
- **QUICKSTART.md**: Step-by-step getting started guide
- **LICENSE**: MIT License
