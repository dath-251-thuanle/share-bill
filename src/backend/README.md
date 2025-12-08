# Share Bill – Backend (Go + Fiber v2 + PostgreSQL)

## Yêu cầu
- Go 1.22 trở lên
- PostgreSQL 16 hoặc 17 (đã cài và đang chạy trên localhost:5432)

## Các bước chạy lần đầu (chỉ làm 1 lần)

1. **Mở project**
   ```bash
   cd D:\University\semester5\DACNPM\dath\share-bill\src\backend

2. **Tạo database trong PostgreSQL**
Mở pgAdmin / DBeaver / psql rồi chạy:
SQLCREATE DATABASE dath2;   -- bạn có thể đổi tên nếu muốn
3. **Tạo bảng**
Mở Query Tool → chạy toàn bộ nội dung file database_create.sql (file này nằm trong thư mục /src/backend)
4. **Tạo file .env**
Tạo file tên .env ngay tại thư mục src/backend/ với nội dung sau:
   --- bash
   envDATABASE_URL=postgres://caotri136:Caotri1306asdfg**@localhost:5432/dath2?sslmode=disable
   APP_PORT=3000
   JWT_SECRET=sharebill-dath2-caotri136-2025-super-secret-key-1306
Lưu ý:
Thay caotri136 và Caotri1306asdfg** bằng username + password PostgreSQL thật của bạn
Thay dath2 bằng tên database bạn vừa tạo ở bước 2

6. **Cài dependencies (chỉ chạy 1 lần)**
   --- bash
   go mod init sharebill-backend        # nếu chưa có go.mod
   go get github.com/gofiber/fiber/v2
   go get github.com/jackc/pgx/v5/pgxpool
   go get github.com/joho/godotenv
   go get golang.org/x/crypto/bcrypt
   go get github.com/golang-jwt/jwt/v5
   go get github.com/go-playground/validator/v10

Chạy server (mỗi lần code)
   --- bash
   cd D:\University\semester5\DACNPM\dath\share-bill\src\backend
   go run cmd/api/main.go

Server sẽ chạy tại: http://localhost:3000
