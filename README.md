# Golang Authentication API Boilerplate

RESTful API authentication boilerplate yang dibangun dengan Go menggunakan framework Gin, ORM GORM, dan JWT untuk autentikasi. API ini menyediakan fitur-fitur authentication lengkap termasuk register, login, reset password via email, dan manajemen profil user.

## ✨ Fitur

- ✅ Registrasi user dengan email verification
- ✅ Login dengan JWT token
- ✅ Logout
- ✅ Forgot password & Reset password via email
- ✅ Change password untuk authenticated user
- ✅ Get & Update user profile
- ✅ Email verification
- ✅ JWT-based authentication
- ✅ Password hashing dengan bcrypt
- ✅ SMTP email service
- ✅ Middleware authentication
- ✅ Request logging middleware
- ✅ CORS enabled
- ✅ Environment-based configuration
- ✅ MySQL/MariaDB database dengan GORM
- ✅ Auto-create database (seperti Eloquent ORM)
- ✅ Struktur project yang terorganisir

## 🛠 Tech Stack

- **Language:** Go 1.26
- **Web Framework:** [Gin](https://github.com/gin-gonic/gin)
- **ORM:** [GORM](https://gorm.io/)
- **Database:** MySQL/MariaDB
- **Authentication:** JWT (JSON Web Tokens)
- **Password Hashing:** bcrypt
- **Email Service:** SMTP
- **Environment Management:** godotenv

## 📁 Struktur Project

```
golang-auth-api-boilerplate/
├── config/              # Konfigurasi aplikasi
│   └── config.go
├── controllers/         # HTTP handlers
│   └── auth_controller.go
├── database/           # Database connection
│   └── database.go
├── middleware/         # Middleware functions
│   ├── auth.go
│   └── logger.go
├── models/             # Data models
│   └── user.go
├── routes/             # Route definitions
│   └── routes.go
├── services/           # Business logic
│   └── email_service.go
├── utils/              # Utility functions
│   ├── helpers.go
│   ├── response.go
│   └── token.go
├── .env.example        # Environment variables template
├── .gitignore
├── go.mod
├── main.go             # Application entry point
└── README.md
```

## 🚀 Getting Started

### Prerequisites

- Go 1.26 atau lebih tinggi
- MySQL 5.7+ atau MariaDB 10.3+
- SMTP server (Gmail, SendGrid, dll) untuk email functionality

### Installation

1. **Clone repository**

```bash
git clone https://github.com/yourusername/golang-auth-api-boilerplate.git
cd golang-auth-api-boilerplate
```

2. **Install dependencies**

```bash
go mod download
```

3. **Setup database**

**Database akan dibuat otomatis!** Anda hanya perlu memastikan MySQL/MariaDB server sudah berjalan.

Jika ingin membuat manual:

```sql
CREATE DATABASE auth_api_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

4. **Configure environment variables**

Copy file `.env.example` ke `.env` dan sesuaikan dengan konfigurasi Anda:

```bash
cp .env.example .env
```

Edit file `.env`:

```env
# Server Configuration
PORT=8080
GIN_MODE=debug

# Database Configuration (MySQL/MariaDB)
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=auth_api_db

# JWT Configuration
JWT_SECRET=your_super_secret_jwt_key_change_this_in_production
JWT_EXPIRATION_HOURS=24

# SMTP Configuration for Email
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your_email@gmail.com
SMTP_PASSWORD=your_app_password
SMTP_FROM=noreply@yourapp.com

# Frontend URL (for password reset links)
FRONTEND_URL=http://localhost:3000
```

**Note untuk Gmail SMTP:**
- Aktifkan 2-Factor Authentication
- Generate App Password dari Google Account settings
- Gunakan App Password sebagai `SMTP_PASSWORD`

5. **Run the application**

```bash
go run main.go
```

Server akan berjalan di `http://localhost:8080`

## 📚 API Documentation

Base URL: `http://localhost:8080/api/v1`

### Health Check

#### GET /health

Cek status server.

**Response:**
```json
{
  "status": "ok",
  "message": "Server is running"
}
```

---

### Authentication Endpoints

#### 1. Register

**POST** `/auth/register`

Registrasi user baru.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123",
  "name": "John Doe"
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "user": {
      "id": 1,
      "email": "user@example.com",
      "name": "John Doe",
      "is_email_verified": false,
      "created_at": "2026-02-16T10:00:00Z",
      "updated_at": "2026-02-16T10:00:00Z"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

**Error Response (409 Conflict):**
```json
{
  "success": false,
  "message": "Email already registered",
  "error": "email_exists"
}
```

---

#### 2. Login

**POST** `/auth/login`

Login user.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "user": {
      "id": 1,
      "email": "user@example.com",
      "name": "John Doe",
      "is_email_verified": true,
      "created_at": "2026-02-16T10:00:00Z",
      "updated_at": "2026-02-16T10:00:00Z"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

**Error Response (401 Unauthorized):**
```json
{
  "success": false,
  "message": "Invalid email or password",
  "error": "invalid_credentials"
}
```

---

#### 3. Forgot Password

**POST** `/auth/forgot-password`

Request password reset link via email.

**Request Body:**
```json
{
  "email": "user@example.com"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "If the email exists, a reset link has been sent"
}
```

---

#### 4. Reset Password

**POST** `/auth/reset-password`

Reset password dengan token yang dikirim via email.

**Request Body:**
```json
{
  "token": "reset_token_from_email",
  "new_password": "newpassword123"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Password reset successfully"
}
```

**Error Response (400 Bad Request):**
```json
{
  "success": false,
  "message": "Invalid or expired reset token",
  "error": "invalid_token"
}
```

---

#### 5. Verify Email

**GET** `/auth/verify-email?token=verification_token`

Verifikasi email address.

**Query Parameters:**
- `token` (required): Verification token dari email

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Email verified successfully"
}
```

---

### Protected Endpoints (Require Authentication)

Semua endpoint berikut memerlukan JWT token di header:

```
Authorization: Bearer <jwt_token>
```

#### 6. Get Profile

**GET** `/user/profile`

Ambil informasi profile user yang sedang login.

**Headers:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Profile retrieved successfully",
  "data": {
    "id": 1,
    "email": "user@example.com",
    "name": "John Doe",
    "is_email_verified": true,
    "created_at": "2026-02-16T10:00:00Z",
    "updated_at": "2026-02-16T10:00:00Z"
  }
}
```

---

#### 7. Update Profile

**PUT** `/user/profile`

Update profile user.

**Headers:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Request Body:**
```json
{
  "name": "John Doe Updated"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Profile updated successfully",
  "data": {
    "id": 1,
    "email": "user@example.com",
    "name": "John Doe Updated",
    "is_email_verified": true,
    "created_at": "2026-02-16T10:00:00Z",
    "updated_at": "2026-02-16T10:00:00Z"
  }
}
```

---

#### 8. Change Password

**POST** `/user/change-password`

Ubah password untuk authenticated user.

**Headers:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Request Body:**
```json
{
  "old_password": "oldpassword123",
  "new_password": "newpassword123"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Password changed successfully"
}
```

**Error Response (401 Unauthorized):**
```json
{
  "success": false,
  "message": "Current password is incorrect",
  "error": "invalid_password"
}
```

---

#### 9. Logout

**POST** `/user/logout`

Logout user (client-side token removal).

**Headers:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Logged out successfully"
}
```

---

## 🔐 Authentication

API ini menggunakan JWT (JSON Web Tokens) untuk authentication. Setelah login atau register, Anda akan menerima token yang harus disertakan di header setiap request ke protected endpoints.

**Format Header:**
```
Authorization: Bearer <your_jwt_token>
```

Token akan expired sesuai dengan konfigurasi `JWT_EXPIRATION_HOURS` di file `.env` (default 24 jam).

## 📧 Email Configuration

### Gmail SMTP Setup

1. Aktifkan 2-Factor Authentication di Google Account
2. Generate App Password:
   - Pergi ke Google Account → Security → App passwords
   - Pilih "Mail" dan device Anda
   - Copy password yang di-generate
3. Gunakan App Password tersebut sebagai `SMTP_PASSWORD` di `.env`

### Other SMTP Providers

Anda bisa menggunakan SMTP provider lain seperti:
- **SendGrid**: smtp.sendgrid.net:587
- **Mailgun**: smtp.mailgun.org:587
- **Amazon SES**: email-smtp.region.amazonaws.com:587

## 🧪 Testing dengan cURL

### Register
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "name": "Test User"
  }'
```

### Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

### Get Profile (with token)
```bash
curl -X GET http://localhost:8080/api/v1/user/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## 🐳 Docker Support (Optional)

Anda bisa menambahkan Docker support dengan membuat `Dockerfile`:

```dockerfile
FROM golang:1.21-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]
```

Dan `docker-compose.yml`:

```yaml
version: '3.8'

services:
  mysql:
    image: mysql:8.0
    environment:
      MYSQL_DATABASE: auth_api_db
      MYSQL_ROOT_PASSWORD: password
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql
    command: --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci

  api:
    build: .
    ports:
      - "8080:8080"
    depends_on:
      - mysql
    environment:
      DB_HOST: mysql
      DB_PORT: 3306
      DB_USER: root
      DB_PASSWORD: password
      DB_NAME: auth_api_db

volumes:
  mysql_data:
```

## 📝 Development Notes

### Database Migrations

GORM akan otomatis membuat/update table schema saat aplikasi start pertama kali (AutoMigrate).

Untuk production, disarankan menggunakan migration tools seperti:
- [golang-migrate](https://github.com/golang-migrate/migrate)
- [goose](https://github.com/pressly/goose)

### Security Best Practices

1. **Jangan commit file `.env`** - Selalu ada di `.gitignore`
2. **Gunakan JWT Secret yang kuat** - Minimal 32 karakter random
3. **Enable HTTPS di production** - Gunakan reverse proxy seperti Nginx
4. **Rate Limiting** - Implementasi rate limiting untuk mencegah brute force
5. **Token Blacklisting** - Untuk production, consider implementing token blacklist untuk logout
6. **Input Validation** - Sudah ada basic validation, bisa ditingkatkan sesuai kebutuhan

### Adding More Features

Struktur project ini mudah untuk di-extend. Beberapa fitur yang bisa ditambahkan:

- Role-based access control (RBAC)
- Refresh tokens
- OAuth2 integration (Google, Facebook, dll)
- Rate limiting middleware
- API key authentication
- User management (admin features)
- Logging dengan structured logger (zerolog, zap)
- Metrics dan monitoring
- Unit tests & integration tests

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the project
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 👤 Author

Your Name - [@yourusername](https://github.com/yourusername)

## 🙏 Acknowledgments

- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [GORM](https://gorm.io/)
- [JWT-Go](https://github.com/golang-jwt/jwt)
- Go Community

---

**Happy Coding! 🚀**
#   g o l a n g - a u t h - a p i - b o i l e r p l a t e 
 
 