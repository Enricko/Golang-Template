# GOLANG TEMPLATE

A modern Golang API template with built-in CLI tools, repository pattern, and comprehensive structure.

## Project Structure

```
.
├── app/
│   ├── controllers/           # HTTP request handlers
│   │   ├── product_controller.go
│   │   └── user_controller.go
│   ├── middleware/           # HTTP middleware
│   │   ├── auth.go
│   │   └── cors.go
│   ├── models/              # Database models
│   │   ├── product.go
│   │   └── user.go
│   └── repositories/        # Repository layer for data access
│       └── user_repository.go
├── cmd/                     # CLI commands
│   ├── key.go              # Key generation commands
│   ├── make.go             # Make commands for scaffolding
│   ├── migrate.go          # Database migration commands
│   ├── root.go             # Root CLI command
│   └── serve.go            # Server command
├── config/                 # Configuration
│   ├── config.go          # Configuration loader
│   └── database.go        # Database configuration
├── database/
│   └── migrations/        # Database migrations
│       └── migration.go   # Migration handler
├── registry/              # Model registry for migrations
│   └── models.go
├── routes/               # HTTP routes
│   └── routes.go        # Route definitions
├── services/            # Business logic layer
├── utils/              # Utility functions
│   ├── jwt.go         # JWT utilities
│   └── password.go    # Password utilities
├── .env               # Environment variables
├── .env.example      # Example environment file
├── go.mod           # Go modules
├── go.sum          # Go modules checksum
└── main.go         # Application entry point
```

## Features

- **Layered Architecture**: Controllers → Services → Repositories → Models
- **CLI Tools**: Built-in commands for code generation and database management
- **JWT Authentication**: Pre-configured JWT authentication
- **Repository Pattern**: Clean separation of data access layer
- **Middleware Support**: Ready-to-use authentication and CORS middleware
- **Environment Management**: Easy configuration through .env files
- **Database Migration**: Built-in migration system
- **Code Generation**: Generate models, controllers, and repositories

## Available Commands

### Start Server
```bash
go run main.go serve
```

### Generate Code
```bash
# Create a new model
go run main.go make:model user

# Create a new controller
go run main.go make:controller user

# Create both model and controller
go run main.go make:controller user -m
```

### Database Migrations
```bash
# Migrate specific table
go run main.go migrate user

# Migrate all tables
go run main.go migrate --all

# Fresh migration (drops all tables and re-migrates)
go run main.go migrate --fresh
```

### Key Generation
```bash
# Generate new key pairs
go run main.go key generate
```

## Configuration

Configuration is managed through environment variables. Copy `.env.example` to `.env`:

```env
APP_NAME=YourApp
APP_ENV=local
APP_PORT=8080

DB_CONNECTION=mysql
DB_HOST=127.0.0.1
DB_PORT=3306
DB_DATABASE=your_database
DB_USERNAME=your_username
DB_PASSWORD=your_password

JWT_SECRET=your-secret-key
```

## Repository Pattern Example

```go
// app/repositories/user_repository.go
type UserRepository struct {
    db *gorm.DB
}

func (r *UserRepository) Find(id uint) (*models.User, error) {
    var user models.User
    result := r.db.First(&user, id)
    return &user, result.Error
}
```

## Controller Example

```go
// app/controllers/user_controller.go
type UserController struct {
    userRepo *repositories.UserRepository
}

func (c *UserController) Show(ctx *gin.Context) {
    id := ctx.Param("id")
    user, err := c.userRepo.Find(id)
    if err != nil {
        ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }
    ctx.JSON(http.StatusOK, user)
}
```

## Middleware Usage

```go
// routes/routes.go
func SetupRoutes(r *gin.Engine) {
    // Public routes
    public := r.Group("/api")
    {
        public.POST("/login", controllers.Login)
    }

    // Protected routes
    protected := r.Group("/api")
    protected.Use(middleware.AuthMiddleware())
    {
        protected.GET("/users", userController.Index)
    }
}
```

## Prerequisites

- Go 1.18 or higher
- MySQL/PostgreSQL/SQLite

## Installation

1. Clone the repository
```bash
git clone [repository-url]
cd [project-name]
```

2. Install dependencies
```bash
go mod download
```

3. Set up environment
```bash
cp .env.example .env
```

4. Run migrations
```bash
go run main.go migrate --all
```

5. Start the server
```bash
go run main.go serve
```

## License

This project is licensed under the MIT License.