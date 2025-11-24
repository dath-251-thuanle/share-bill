package main

import (
    "fmt"
    "log"
    "net/http"
    "os"

    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
)

func main() {
    // Load .env file
    if err := godotenv.Load("../../.env"); err != nil {
        log.Println("Warning: No .env file found, using environment variables")
    }

    // Get port from environment
    port := os.Getenv("APP_PORT")
    if port == "" {
        port = "8080"
    }

    // Initialize router
    router := gin.Default()

    // Health check endpoint
    router.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "status":  "healthy",
            "message": "Share Bill Backend API is running",
            "version": "0.1.0",
        })
    })

    // Ping endpoint
    router.GET("/api/v1/ping", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "message": "pong",
        })
    })

    // Placeholder routes
    v1 := router.Group("/api/v1")
    {
        v1.GET("/events", func(c *gin.Context) {
            c.JSON(http.StatusOK, gin.H{
                "message": "List events - Coming soon",
            })
        })
        
        v1.GET("/members", func(c *gin.Context) {
            c.JSON(http.StatusOK, gin.H{
                "message": "List members - Coming soon",
            })
        })
        
        v1.GET("/transactions", func(c *gin.Context) {
            c.JSON(http.StatusOK, gin.H{
                "message": "List transactions - Coming soon",
            })
        })
    }

    // Start server
    addr := fmt.Sprintf(":%s", port)
    log.Printf("🚀 Server starting on http://localhost%s", addr)
    log.Printf("📚 Health check: http://localhost%s/health", addr)
    
    if err := router.Run(addr); err != nil {
        log.Fatalf("❌ Failed to start server: %v", err)
    }
}
