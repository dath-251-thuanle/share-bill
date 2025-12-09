// cmd/api/main.go
package main

import (
	"log"
	"os"

	"sharebill-backend/internal/app"
	"sharebill-backend/internal/config"
	"sharebill-backend/internal/handlers"
	"sharebill-backend/internal/repositories"
	"sharebill-backend/internal/services"
)

func main() {
	// 1. Kết nối DB
	db := config.InitDB()
	defer db.Close()

	// 1.1 khởi tạo chức năng liên quan đến user
	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo)
	authHandler := handlers.NewAuthHandler(authService)

	// 2. Khởi tạo Fiber app
	fiberApp := app.New()

	// 2.1 Đăng ký route cho auth
	api := fiberApp.Group("/api")
	api.Post("/auth/register", authHandler.Register)
	api.Post("/auth/login", authHandler.Login)

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