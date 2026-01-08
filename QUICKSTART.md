# Quick Start Guide

## 1. Get NVIDIA API Keys

1. Go to https://build.nvidia.com/
2. Sign up for a free account
3. Navigate to any model page (e.g., minimax-m2)
4. Click "Get API Key"
5. Copy your API key (starts with "nvapi-")
6. Repeat with multiple accounts to get more API keys

## 2. Configure ProxyPal

```bash
# Copy example config
cp config.example.yaml config.yaml

# Edit config.yaml and add your API keys
# Replace the example keys with your real ones:
nvidia:
  api_keys:
    - "nvapi-YOUR-FIRST-KEY-HERE"
    - "nvapi-YOUR-SECOND-KEY-HERE"
    - "nvapi-YOUR-THIRD-KEY-HERE"
```

## 3. Run ProxyPal

### Option A: Using Docker (Recommended)

```bash
# Build and start
docker-compose up -d

# Check logs
docker-compose logs -f

# Test
curl http://localhost:8080/health
```

### Option B: Using Binary

```bash
# Build
./build.sh

# Run
./proxypal
```

### Option C: Using Make

```bash
# Build and run
make dev

# Or step by step
make build
make run
```

## 4. Test Your Setup

### Test with curl:

```bash
# Health check
curl http://localhost:8080/health

# Get statistics
curl http://localhost:8080/stats

# List models
curl http://localhost:8080/v1/models

# Chat completion
curl -X POST http://localhost:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "minimaxai/minimax-m2",
    "messages": [{"role": "user", "content": "Hello!"}],
    "max_tokens": 100
  }'
```

### Test with Python:

```bash
# Install OpenAI library
pip install openai

# Run example
python example.py
```

## 5. Use in Your Application

### Python Example:

```python
from openai import OpenAI

client = OpenAI(
    base_url="http://localhost:8080/v1",
    api_key="any-value-works"
)

response = client.chat.completions.create(
    model="minimaxai/minimax-m2",
    messages=[{"role": "user", "content": "Your prompt here"}],
    stream=True
)

for chunk in response:
    if chunk.choices[0].delta.content:
        print(chunk.choices[0].delta.content, end="")
```

### Node.js Example:

```javascript
import OpenAI from 'openai';

const client = new OpenAI({
  baseURL: 'http://localhost:8080/v1',
  apiKey: 'any-value-works'
});

const stream = await client.chat.completions.create({
  model: 'minimaxai/minimax-m2',
  messages: [{ role: 'user', content: 'Your prompt here' }],
  stream: true,
});

for await (const chunk of stream) {
  process.stdout.write(chunk.choices[0]?.delta?.content || '');
}
```

## 6. Monitor Usage

```bash
# Check statistics
curl http://localhost:8080/stats | jq

# View logs (if using Docker)
docker-compose logs -f
```

## Troubleshooting

**Problem: Server won't start**
- Check if port 8080 is already in use: `lsof -i :8080`
- Verify config.yaml exists and is valid
- Check logs for error messages

**Problem: Rate limit errors**
- Check `/stats` endpoint to see token availability
- Add more API keys to increase capacity
- Wait for tokens to refill (40 per minute per key)

**Problem: Invalid API key errors**
- Verify your NVIDIA API keys are correct
- Make sure keys start with "nvapi-"
- Check keys haven't expired on build.nvidia.com

## Advanced Configuration

### Change Port:

```yaml
server:
  port: 3000  # Use different port
```

### Adjust Rate Limit:

```yaml
nvidia:
  rate_limit: 40  # Requests per minute per key
```

### Enable Debug Logging:

```yaml
logging:
  level: "debug"
  enable_request_log: true
```

## Next Steps

- Read the full [README.md](README.md) for more details
- Check available NVIDIA models at https://build.nvidia.com/
- Integrate ProxyPal into your application
- Deploy to production with Docker

## Support

- Issues: https://github.com/luongcoder/proxypal-nvidia-load-balance/issues
- Documentation: [README.md](README.md)
