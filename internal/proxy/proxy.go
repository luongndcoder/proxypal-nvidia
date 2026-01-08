package proxy

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/luongcoder/proxypal-nvidia-load-balance/internal/balancer"
	"github.com/luongcoder/proxypal-nvidia-load-balance/internal/config"
)

// ProxyServer handles incoming requests and proxies them to NVIDIA API
type ProxyServer struct {
	loadBalancer *balancer.LoadBalancer
	config       *config.Config
	httpClient   *http.Client
}

// NewProxyServer creates a new proxy server
func NewProxyServer(cfg *config.Config, lb *balancer.LoadBalancer) *ProxyServer {
	return &ProxyServer{
		loadBalancer: lb,
		config:       cfg,
		httpClient: &http.Client{
			Timeout: time.Duration(cfg.NVIDIA.Timeout) * time.Second,
		},
	}
}

// SetupRoutes configures the Gin router with proxy endpoints
func (ps *ProxyServer) SetupRoutes(router *gin.Engine) {
	// OpenAI-compatible endpoints
	v1 := router.Group("/v1")
	{
		v1.POST("/chat/completions", ps.handleChatCompletions)
		v1.GET("/models", ps.handleListModels)
	}

	// Health check and stats endpoints
	router.GET("/health", ps.handleHealth)
	router.GET("/stats", ps.handleStats)
}

// handleChatCompletions proxies chat completion requests to NVIDIA API
func (ps *ProxyServer) handleChatCompletions(c *gin.Context) {
	// Read request body
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read request body"})
		return
	}

	// Parse request to check if it's streaming
	var reqBody map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON request"})
		return
	}

	isStreaming := false
	if stream, ok := reqBody["stream"].(bool); ok {
		isStreaming = stream
	}

	// Get API key from load balancer
	apiKey, err := ps.loadBalancer.GetKeyWithRetry(ps.config.NVIDIA.Retry.MaxRetries)
	if err != nil {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error": map[string]interface{}{
				"message": err.Error(),
				"type":    "rate_limit_error",
				"code":    "rate_limit_exceeded",
			},
		})
		return
	}

	// Log request if enabled
	if ps.config.Logging.EnableRequestLog {
		model := "unknown"
		if m, ok := reqBody["model"].(string); ok {
			model = m
		}
		fmt.Printf("[%s] Request to model: %s, streaming: %v, key: %s\n",
			time.Now().Format("2006-01-02 15:04:05"),
			model,
			isStreaming,
			balancer.MaskAPIKey(apiKey.Key))
	}

	// Create request to NVIDIA API
	url := ps.config.NVIDIA.BaseURL + "/chat/completions"
	req, err := http.NewRequest("POST", url, bytes.NewReader(bodyBytes))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
		return
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey.Key)

	// Forward other headers from original request
	for key, values := range c.Request.Header {
		if key != "Authorization" && key != "Host" {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}
	}

	// Execute request
	resp, err := ps.httpClient.Do(req)
	if err != nil {
		ps.loadBalancer.MarkKeyError(apiKey)
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to contact NVIDIA API"})
		return
	}
	defer resp.Body.Close()

	// Handle rate limit errors
	if resp.StatusCode == http.StatusTooManyRequests {
		ps.loadBalancer.MarkKeyError(apiKey)

		// Try with a different key if auto-failover is enabled
		if ps.config.NVIDIA.Retry.AutoFailover {
			newKey, err := ps.loadBalancer.GetKeyWithRetry(ps.config.NVIDIA.Retry.MaxRetries)
			if err == nil {
				// Retry with new key
				req.Header.Set("Authorization", "Bearer "+newKey.Key)
				resp, err = ps.httpClient.Do(req)
				if err != nil {
					c.JSON(http.StatusBadGateway, gin.H{"error": "failed to contact NVIDIA API"})
					return
				}
				defer resp.Body.Close()
			}
		}
	}

	// Copy response headers
	for key, values := range resp.Header {
		for _, value := range values {
			c.Header(key, value)
		}
	}

	// Handle streaming response
	if isStreaming {
		ps.handleStreamingResponse(c, resp)
		return
	}

	// Handle non-streaming response
	c.Status(resp.StatusCode)
	io.Copy(c.Writer, resp.Body)
}

// handleStreamingResponse handles server-sent events streaming
func (ps *ProxyServer) handleStreamingResponse(c *gin.Context, resp *http.Response) {
	c.Status(resp.StatusCode)
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	// Flush headers
	c.Writer.Flush()

	// Stream the response
	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err != io.EOF {
				fmt.Printf("Error reading stream: %v\n", err)
			}
			break
		}

		// Write line to client
		c.Writer.Write(line)
		c.Writer.Flush()
	}
}

// handleListModels returns available models
func (ps *ProxyServer) handleListModels(c *gin.Context) {
	// Get API key
	apiKey, err := ps.loadBalancer.GetNextKey()
	if err != nil {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": err.Error()})
		return
	}

	// Create request to NVIDIA API
	url := ps.config.NVIDIA.BaseURL + "/models"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
		return
	}

	req.Header.Set("Authorization", "Bearer "+apiKey.Key)

	// Execute request
	resp, err := ps.httpClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to contact NVIDIA API"})
		return
	}
	defer resp.Body.Close()

	// Forward response
	c.Status(resp.StatusCode)
	io.Copy(c.Writer, resp.Body)
}

// handleHealth returns health status
func (ps *ProxyServer) handleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"time":   time.Now().Format(time.RFC3339),
	})
}

// handleStats returns load balancer statistics
func (ps *ProxyServer) handleStats(c *gin.Context) {
	stats := ps.loadBalancer.GetStats()

	c.JSON(http.StatusOK, gin.H{
		"keys":      len(stats),
		"stats":     stats,
		"timestamp": time.Now().Format(time.RFC3339),
	})
}
