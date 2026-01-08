# ProxyPal NVIDIA Load Balancer

A high-performance load balancer for NVIDIA API keys, written in Golang. This tool helps you maximize free tier usage by distributing requests across multiple NVIDIA API accounts with intelligent rate limiting.

## Features

- **OpenAI-Compatible API**: Drop-in replacement for OpenAI client libraries
- **Round-Robin Load Balancing**: Automatically distributes requests across multiple API keys
- **Rate Limiting**: Respects NVIDIA's 40 requests/minute limit per API key
- **Auto-Failover**: Automatically switches to another key if one is rate-limited
- **Streaming Support**: Full support for server-sent events (SSE) streaming
- **Statistics Dashboard**: Monitor usage and performance of each API key
- **Docker Support**: Easy deployment with Docker and docker-compose
- **Zero Dependencies**: Compiled to a single binary

## Quick Start

### Using Docker (Recommended)

1. **Clone the repository**
   ```bash
   git clone https://github.com/luongcoder/proxypal-nvidia-load-balance.git
   cd proxypal-nvidia-load-balance
   ```

2. **Create configuration file**
   ```bash
   cp config.example.yaml config.yaml
   ```

3. **Edit config.yaml and add your NVIDIA API keys**
   ```yaml
   nvidia:
     api_keys:
       - "nvapi-your-first-key-here"
       - "nvapi-your-second-key-here"
       - "nvapi-your-third-key-here"
   ```

4. **Start with Docker Compose**
   ```bash
   docker-compose up -d
   ```

5. **Test the service**
   ```bash
   curl http://localhost:8080/health
   ```

### Using Binary

1. **Install Go 1.21 or higher**

2. **Build the application**
   ```bash
   chmod +x build.sh
   ./build.sh
   ```

3. **Configure your API keys**
   ```bash
   cp config.example.yaml config.yaml
   # Edit config.yaml with your favorite editor
   ```

4. **Run the application**
   ```bash
   ./proxypal
   ```

## Usage

### With Python (OpenAI Library)

```python
from openai import OpenAI

client = OpenAI(
    base_url="http://localhost:8080/v1",
    api_key="dummy-key"  # Any value works, will be replaced by load balancer
)

completion = client.chat.completions.create(
    model="minimaxai/minimax-m2",
    messages=[{"role": "user", "content": "Hello, how are you?"}],
    temperature=1,
    top_p=0.95,
    max_tokens=8192,
    stream=True
)

for chunk in completion:
    if chunk.choices[0].delta.content is not None:
        print(chunk.choices[0].delta.content, end="")
```

### With cURL

**Non-streaming request:**
```bash
curl -X POST http://localhost:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "minimaxai/minimax-m2",
    "messages": [{"role": "user", "content": "Hello!"}],
    "temperature": 1,
    "max_tokens": 8192
  }'
```

**Streaming request:**
```bash
curl -X POST http://localhost:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "minimaxai/minimax-m2",
    "messages": [{"role": "user", "content": "Hello!"}],
    "temperature": 1,
    "max_tokens": 8192,
    "stream": true
  }'
```

## Configuration

The `config.yaml` file supports the following options:

```yaml
server:
  port: 8080              # Server port
  host: "0.0.0.0"         # Bind address

nvidia:
  base_url: "https://integrate.api.nvidia.com/v1"
  rate_limit: 40          # Requests per minute per key
  api_keys:               # List your API keys here
    - "nvapi-key-1"
    - "nvapi-key-2"
    - "nvapi-key-3"
  timeout: 300            # Request timeout in seconds

  retry:
    max_retries: 3        # Max retries when rate limited
    auto_failover: true   # Auto switch to another key

logging:
  level: "info"           # Log level: debug, info, warn, error
  enable_request_log: true
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/v1/chat/completions` | OpenAI-compatible chat completions |
| GET | `/v1/models` | List available models |
| GET | `/health` | Health check endpoint |
| GET | `/stats` | Load balancer statistics |

### View Statistics

```bash
curl http://localhost:8080/stats
```

Example response:
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
    },
    {
      "KeyPrefix": "nvapi-...xyz2",
      "RequestCount": 145,
      "ErrorCount": 0,
      "AvailableTokens": 40,
      "LastUsed": "2024-01-08T10:29:55Z"
    }
  ],
  "timestamp": "2024-01-08T10:30:05Z"
}
```

## How It Works

1. **Round-Robin Selection**: Requests are distributed evenly across all API keys
2. **Token Bucket Rate Limiting**: Each key has 40 tokens (requests) per minute
3. **Automatic Refill**: Tokens refill continuously based on elapsed time
4. **Smart Failover**: If a key is rate-limited, automatically tries the next available key
5. **Transparent Proxying**: All requests are forwarded to NVIDIA API with the selected key

## Performance

With 3 API keys, you can achieve:
- **~120 requests/minute** total capacity
- **2 requests/second** sustained throughput
- Automatic load distribution
- High availability with failover

## Development

### Project Structure

```
.
├── cmd/
│   └── proxypal/         # Main application
│       └── main.go
├── internal/
│   ├── balancer/         # Load balancing logic
│   │   ├── loadbalancer.go
│   │   └── ratelimiter.go
│   ├── config/           # Configuration management
│   │   └── config.go
│   └── proxy/            # HTTP proxy server
│       └── proxy.go
├── config.example.yaml   # Example configuration
├── Dockerfile
├── docker-compose.yml
├── build.sh
└── README.md
```

### Building from Source

```bash
# Clone repository
git clone https://github.com/luongcoder/proxypal-nvidia-load-balance.git
cd proxypal-nvidia-load-balance

# Download dependencies
go mod download

# Build
go build -o proxypal ./cmd/proxypal

# Run
./proxypal
```

### Running Tests

```bash
go test ./...
```

## Docker Deployment

### Build Image

```bash
docker build -t proxypal-nvidia:latest .
```

### Run Container

```bash
docker run -d \
  -p 8080:8080 \
  -v $(pwd)/config.yaml:/root/config.yaml \
  -e CONFIG_PATH=/root/config.yaml \
  --name proxypal \
  proxypal-nvidia:latest
```

### Using Docker Compose

```bash
# Start
docker-compose up -d

# View logs
docker-compose logs -f

# Stop
docker-compose down
```

## FAQ

**Q: Is this legal?**
A: Yes! You're using NVIDIA's free tier as intended. You're just managing multiple accounts efficiently.

**Q: How many API keys should I use?**
A: Start with 3-5 keys. More keys = higher throughput, but manage them responsibly.

**Q: What happens if all keys are rate-limited?**
A: The server returns a 429 status code. The client should implement exponential backoff.

**Q: Can I use this in production?**
A: Yes, but monitor your usage and ensure you comply with NVIDIA's terms of service.

**Q: Does it support all NVIDIA models?**
A: Yes! It's a transparent proxy - any model available through NVIDIA's API will work.

## Troubleshooting

**Server won't start:**
- Check if port 8080 is already in use
- Verify your config.yaml is valid YAML
- Ensure at least one API key is configured

**All requests fail with rate limit errors:**
- Check if your API keys are valid
- Verify the rate_limit setting (default: 40 req/min)
- Check `/stats` endpoint to see token availability

**Streaming doesn't work:**
- Ensure you're using HTTP/1.1 (not HTTP/2)
- Set `stream: true` in your request body
- Check client library supports SSE

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Disclaimer

This tool is for educational purposes and managing your own NVIDIA API keys efficiently. Users are responsible for complying with NVIDIA's Terms of Service and API usage policies.

## Support

- Create an issue: [GitHub Issues](https://github.com/luongcoder/proxypal-nvidia-load-balance/issues)
- Star the project if you find it useful!

---

**Happy Load Balancing! 🚀**
