// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter "Bearer <token>"

package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "ecommerce/internal/db"
    httpHandler "ecommerce/internal/handler/http"
    "ecommerce/internal/repository/postgres"
    authService "ecommerce/internal/service/auth"
    userService "ecommerce/internal/service/user"
    productService "ecommerce/internal/service/product"

    _ "github.com/lib/pq"
    _ "ecommerce/docs"

    "github.com/joho/godotenv"
)

// @title Ecommerce API
// @version 1.0
// @host localhost:8080
// @BasePath /api/v1

func main() {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using environment variables")
    }

    dbHost := getEnv("DB_HOST", "localhost")
    dbPort := getEnv("DB_PORT", "5432")
    dbUser := getEnv("DB_USER", "postgres")
    dbPassword := getEnv("DB_PASSWORD", "postgres")
    dbName := getEnv("DB_NAME", "ecommerce_db")
    dbSSLMode := getEnv("DB_SSLMODE", "disable")
    serverPort := getEnv("SERVER_PORT", "8080")
    jwtSecret := getEnv("JWT_SECRET", "your-secret-key-change-in-production")

    dbConfig := db.Config{
        Host:     dbHost,
        Port:     dbPort,
        User:     dbUser,
        Password: dbPassword,
        DBName:   dbName,
        SSLMode:  dbSSLMode,
    }

    database, err := db.NewDB(dbConfig)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer database.Close()

    log.Println("Database connection established")

    // -------------------------
    // REPOSITORIES
    // -------------------------
    userRepo := postgres.NewUserRepository(database)
    productRepo := postgres.NewProductRepository(database.DB)

    // -------------------------
    // SERVICES
    // -------------------------
    authSvc := authService.NewService(userRepo, jwtSecret)
    userSvc := userService.NewService(userRepo)
    productSvc := productService.NewService(productRepo)

    // -------------------------
    // ROUTER
    // -------------------------
    router := httpHandler.NewRouter(httpHandler.RouterConfig{
        AuthService:    authSvc,
        UserService:    userSvc,
        ProductService: productSvc, 
    })

    server := &http.Server{
        Addr:         ":" + serverPort,
        Handler:      router,
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 15 * time.Second,
        IdleTimeout:  60 * time.Second,
    }

    go func() {
        log.Printf("Server starting on http://localhost:%s", serverPort)
        log.Printf("API documentation: http://localhost:%s/swagger/index.html", serverPort)
        log.Printf("Health check: http://localhost:%s/health", serverPort)

        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Server failed to start: %v", err)
        }
    }()

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    log.Println("Shutting down server...")

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := server.Shutdown(ctx); err != nil {
        log.Fatalf("Server forced to shutdown: %v", err)
    }

    log.Println("Server stopped gracefully")
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
