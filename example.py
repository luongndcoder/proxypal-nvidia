#!/usr/bin/env python3
"""
Example usage of ProxyPal NVIDIA Load Balancer with OpenAI Python client
"""

from openai import OpenAI

# Create client pointing to ProxyPal
client = OpenAI(
    base_url="http://localhost:8080/v1",
    api_key="dummy-key"  # Can be any value, will be replaced by load balancer
)

def test_non_streaming():
    """Test non-streaming chat completion"""
    print("\n" + "="*60)
    print("Testing Non-Streaming Request")
    print("="*60)

    response = client.chat.completions.create(
        model="minimaxai/minimax-m2",
        messages=[
            {"role": "user", "content": "Write a haiku about artificial intelligence"}
        ],
        temperature=1,
        max_tokens=100
    )

    print(f"\nResponse: {response.choices[0].message.content}")
    print(f"Tokens used: {response.usage.total_tokens}")

def test_streaming():
    """Test streaming chat completion"""
    print("\n" + "="*60)
    print("Testing Streaming Request")
    print("="*60 + "\n")

    stream = client.chat.completions.create(
        model="minimaxai/minimax-m2",
        messages=[
            {"role": "user", "content": "Tell me a short joke about programming"}
        ],
        temperature=1,
        max_tokens=100,
        stream=True
    )

    print("Response: ", end="", flush=True)
    for chunk in stream:
        if chunk.choices[0].delta.content is not None:
            print(chunk.choices[0].delta.content, end="", flush=True)
    print("\n")

def test_list_models():
    """Test listing available models"""
    print("\n" + "="*60)
    print("Testing List Models")
    print("="*60)

    models = client.models.list()
    print(f"\nAvailable models: {len(models.data)}")
    for model in models.data[:5]:  # Show first 5 models
        print(f"  - {model.id}")

if __name__ == "__main__":
    try:
        print("\n" + "="*60)
        print("ProxyPal NVIDIA Load Balancer - Example Usage")
        print("="*60)
        print("\nMake sure ProxyPal is running on http://localhost:8080")
        print("Start it with: ./proxypal")

        # Run tests
        test_list_models()
        test_non_streaming()
        test_streaming()

        print("\n" + "="*60)
        print("All tests completed successfully!")
        print("="*60 + "\n")

    except Exception as e:
        print(f"\nError: {e}")
        print("\nMake sure ProxyPal is running and your API keys are configured correctly.")
