package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/luongndcoder/proxypal-nvidia/internal/balancer"
	"github.com/luongndcoder/proxypal-nvidia/internal/config"
	"github.com/luongndcoder/proxypal-nvidia/internal/proxy"
)

func main() {
	// Load configuration
	configPath := getConfigPath()
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Set log level
	if cfg.Logging.Level != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize load balancer
	lb := balancer.NewLoadBalancer(&cfg.NVIDIA)

	// Create proxy server
	proxyServer := proxy.NewProxyServer(cfg, lb)

	// Setup Gin router
	router := gin.Default()

	// Add CORS middleware
	router.Use(corsMiddleware())

	// Setup routes
	proxyServer.SetupRoutes(router)

	// Print startup info
	printStartupInfo(cfg, lb)

	// Start server
	addr := cfg.GetAddress()
	log.Printf("Starting ProxyPal NVIDIA Load Balancer on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// getConfigPath returns the configuration file path
func getConfigPath() string {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.yaml"
	}
	return configPath
}

// corsMiddleware adds CORS headers
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// printStartupInfo prints startup information
func printStartupInfo(cfg *config.Config, lb *balancer.LoadBalancer) {
	stats := lb.GetStats()

	banner := "============================================================"
	fmt.Println("\n" + banner)
	fmt.Println("  ProxyPal NVIDIA Load Balancer")
	fmt.Println(banner)
	fmt.Printf("  Server Address: http://%s\n", cfg.GetAddress())
	fmt.Printf("  API Keys: %d\n", len(stats))
	fmt.Printf("  Rate Limit: %d requests/minute per key\n", cfg.NVIDIA.RateLimit)
	fmt.Printf("  Total Capacity: ~%d requests/minute\n", cfg.NVIDIA.RateLimit*len(stats))
	fmt.Println(banner)
	fmt.Println("\n  Endpoints:")
	fmt.Printf("    POST   /v1/chat/completions   - OpenAI-compatible chat completions\n")
	fmt.Printf("    GET    /v1/models             - List available models\n")
	fmt.Printf("    GET    /health                - Health check\n")
	fmt.Printf("    GET    /stats                 - Load balancer statistics\n")
	fmt.Println("\n" + banner)
	fmt.Println()
}
