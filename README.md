# Evermos Rakamin Backend

Backend API untuk project Evermos Rakamin dengan teknologi Go + Fiber, MySQL, dan JWT authentication.

## 🚀 Tech Stack

- **Framework:** [Fiber](https://gofiber.io/) (Fast HTTP framework for Go)
- **Database:** MySQL with [GORM](https://gorm.io/) ORM
- **Authentication:** JWT (JSON Web Tokens)
- **Validation:** [go-playground/validator](https://github.com/go-playground/validator)
- **Password Hashing:** bcrypt via golang.org/x/crypto
- **External API:** API Wilayah Indonesia integration

## 📋 Features

- ✅ User Authentication (Register, Login, Profile) dengan JWT
- ✅ Category CRUD dengan pagination
- ✅ Product CRUD dengan search, pagination, dan category relationship
- ✅ API Wilayah Indonesia integration (Provinces, Regencies, Districts)
- ✅ Database migrations dengan GORM AutoMigrate
- ✅ Seed script untuk sample data
- ✅ Health check endpoints (service & database)

## 📁 Project Structure

```
.
├── cmd/
│   └── seed/              # Seed script command
├── config/                 # Configuration loader
├── constants/             # App constants (roles, etc)
├── domain/                 # Domain models & DTOs
│   ├── dto/               # Data Transfer Objects
│   ├── user.go            # User model
│   ├── category.go        # Category model
│   └── product.go         # Product model
├── handlers/               # HTTP handlers (Fiber routes)
├── helpers/                # Helper functions (password, JWT)
├── middleware/             # Middleware (auth, logging)
├── repositories/           # Database access layer
│   ├── db.go              # Database connection
│   ├── errors.go          # Shared errors
│   ├── user_repository.go
│   ├── category_repository.go
│   └── product_repository.go
├── services/               # Business logic layer
│   ├── auth_service.go
│   ├── category_service.go
│   ├── product_service.go
│   └── wilayah_service.go
├── utils/                  # Utilities (response, pagination)
├── .env.example           # Environment variables template
├── go.mod                  # Go modules
├── go.sum                  # Go dependencies checksum
└── main.go                # Application entry point
```

## 🔧 Installation & Setup

### Prerequisites

- Go 1.21+ installed
- MySQL 8.0+ installed and running
- Git

### Steps

1. **Clone repository:**
```bash
git clone https://github.com/rdsarjito/evermos-rakamin.git
cd evermos-rakamin
```

2. **Install dependencies:**
```bash
go mod download
```

3. **Configure environment:**
```bash
cp .env.example .env
```

Edit `.env` file dengan konfigurasi database dan secret key:
```env
APP_HOST=localhost
APP_PORT=8080
SECRET_KEY=your-secret-key-min-32-chars-replace-this

DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=root
DB_PASSWORD=yourpassword
DB_NAME=evermos_clone

API_LOCATION=https://www.emsifa.com/api-wilayah-indonesia/api
```

4. **Create database:**
```sql
CREATE DATABASE evermos_clone;
```

5. **Run seed (optional):**
```bash
go run cmd/seed/main.go
```

6. **Run application:**
```bash
go run main.go
```

Server akan berjalan di `http://localhost:8080`

## 📚 API Endpoints

### Health Check
- `GET /health` - Service health check
- `GET /health/db` - Database health check

### Authentication
- `POST /auth/register` - Register new user
  ```json
  {
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123",
    "role": "user" // optional: admin, user, seller
  }
  ```
- `POST /auth/login` - Login user
  ```json
  {
    "email": "john@example.com",
    "password": "password123"
  }
  ```
- `GET /auth/profile` - Get user profile (requires JWT token)

### Categories
- `GET /categories?page=1&limit=10` - List categories (public)
- `GET /categories/:id` - Get category by ID (public)
- `POST /categories` - Create category (requires JWT)
  ```json
  {
    "name": "Electronics",
    "description": "Electronic devices"
  }
  ```
- `PUT /categories/:id` - Update category (requires JWT)
- `DELETE /categories/:id` - Delete category (requires JWT)

### Products
- `GET /products?page=1&limit=10&search=laptop&category_id=1` - List products (public)
  - Query params: `page`, `limit`, `search`, `category_id`
- `GET /products/:id` - Get product by ID (public)
- `POST /products` - Create product (requires JWT)
  ```json
  {
    "name": "Laptop",
    "description": "High performance laptop",
    "price": 8999.99,
    "stock": 10,
    "category_id": 1
  }
  ```
- `PUT /products/:id` - Update product (requires JWT)
- `DELETE /products/:id` - Delete product (requires JWT)

### Locations (Wilayah Indonesia)
- `GET /locations/provinces` - Get all provinces (public)
- `GET /locations/regencies?province_id=32` - Get regencies by province ID (public)
- `GET /locations/districts?regency_id=3273` - Get districts by regency ID (public)

## 🔐 Authentication

API menggunakan JWT Bearer token untuk authentication. Setelah login, gunakan token di header:

```
Authorization: Bearer <your-jwt-token>
```

Token berlaku selama 24 jam.

## 📝 Example API Calls

### Register & Login
```bash
# Register
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Test User","email":"test@example.com","password":"password123"}'

# Login
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# Use token for protected routes
TOKEN="your-jwt-token-here"
curl -X GET http://localhost:8080/auth/profile \
  -H "Authorization: Bearer $TOKEN"
```

### Category Operations
```bash
# List categories
curl http://localhost:8080/categories?page=1&limit=10

# Create category
curl -X POST http://localhost:8080/categories \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Electronics","description":"Electronic devices"}'
```

### Product Operations
```bash
# Search products
curl "http://localhost:8080/products?search=laptop&category_id=1"

# Create product
curl -X POST http://localhost:8080/products \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Laptop","description":"High performance","price":8999.99,"stock":10,"category_id":1}'
```

## 🗄️ Database

Database menggunakan MySQL dengan GORM ORM. Tables akan otomatis dibuat saat pertama kali run aplikasi melalui AutoMigrate.

### Tables:
- `users` - User accounts
- `categories` - Product categories
- `products` - Products with category relationship

### Seed Data:
Seed script akan membuat sample data:
- 3 users (admin, user, seller) dengan password: `password123`
- 4 categories (Electronics, Clothing, Food & Beverages, Books)
- 6 sample products

## 🔒 Security Features

- Password hashing dengan bcrypt
- JWT token authentication
- Protected routes dengan middleware
- Input validation dengan validator
- SQL injection protection via GORM

## 🚧 Next Steps (Optional Enhancements)

- [ ] Global error handler middleware
- [ ] Structured logging dengan request ID
- [ ] Swagger/OpenAPI documentation
- [ ] Docker & docker-compose setup
- [ ] CI/CD dengan GitHub Actions
- [ ] Unit tests & integration tests
- [ ] Rate limiting
- [ ] CORS configuration
- [ ] Request timeout middleware

## 📄 License

This project is part of Rakamin Academy activities.

## 👤 Author

Built as part of Evermos Rakamin project.

## 🙏 Acknowledgments

- [Fiber Framework](https://gofiber.io/)
- [GORM](https://gorm.io/)
- [API Wilayah Indonesia](https://github.com/emsifa/api-wilayah-indonesia)

