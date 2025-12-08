// cmd/api/main.go
package main

import (
	"log"
	"os"

	"sharebill-backend/internal/app"
	"sharebill-backend/internal/config"
)

func main() {
	// 1. Kết nối DB
	db := config.InitDB()
	defer db.Close()

	// 2. Khởi tạo Fiber app
	fiberApp := app.New()

	// 3. Lấy port
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}

	// 4. Chạy server
	log.Printf("Server running on http://localhost:%s", port)
	log.Printf("Health check: http://localhost:%s/health", port)
	log.Fatal(fiberApp.Listen(":" + port))
}