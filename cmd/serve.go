// cmd/serve.go
package cmd

import (
    "log"
    "golang-template/config"
    "golang-template/database/migrations"
    "golang-template/app/middleware"
    "golang-template/routes"
    "github.com/gin-gonic/gin"
    "github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
    Use:   "serve",
    Short: "Start the application server",
    Run: func(cmd *cobra.Command, args []string) {
        startServer()
    },
}

func init() {
    rootCmd.AddCommand(serveCmd)
}

func startServer() {
    // Load configuration
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatal("Failed to load configuration:", err)
    }

    // Initialize database
    db, err := config.InitDB(cfg)
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    // Run migrations
    if err := migrations.Migrate(db); err != nil {
        log.Fatal("Failed to run migrations:", err)
    }

    // Initialize Gin
    r := gin.Default()

    // Apply middleware
    r.Use(middleware.CorsMiddleware())

    // Setup routes
    routes.SetupRouter(r)

    // Start server
    log.Printf("Server starting on port %s", cfg.AppPort)
    if err := r.Run(":" + cfg.AppPort); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}