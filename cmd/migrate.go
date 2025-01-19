// cmd/migrate.go
package cmd

import (
	"fmt"
	"golang-template/app/models"
	"golang-template/config"
	"golang-template/database/migrations"
	"golang-template/registry"
	"strings"

	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var (
    freshFlag bool
    allFlag bool
)

var migrateCmd = &cobra.Command{
    Use:   "migrate [table_name]",
    Short: "Run database migrations",
    Long: `Run database migrations. You can migrate a specific table or all tables.
Examples:
  migrate user          # Migrate user table
  migrate product      # Migrate product table
  migrate --all        # Migrate all tables
  migrate --fresh      # Drop all tables and re-migrate`,
    Run: func(cmd *cobra.Command, args []string) {
        if freshFlag {
            runFreshMigrations()
            return
        }
        
        if allFlag || len(args) == 0 {
            runAllMigrations()
            return
        }
        
        runSpecificMigration(args[0])
    },
}

func init() {
    migrateCmd.Flags().BoolVarP(&freshFlag, "fresh", "f", false, "Drop all tables and re-migrate")
    migrateCmd.Flags().BoolVarP(&allFlag, "all", "a", false, "Migrate all tables")
    rootCmd.AddCommand(migrateCmd)
}

func runAllMigrations() {
    db := initDB()
    if db == nil {
        return
    }

    fmt.Println("Running all migrations...")
    
    if err := migrations.Migrate(db); err != nil {
        fmt.Printf("Failed to run migrations: %v\n", err)
        return
    }

    fmt.Println("All migrations completed successfully")
}

func runFreshMigrations() {
    db := initDB()
    if db == nil {
        return
    }

    fmt.Println("Dropping all tables...")
    
    // Get all registered models
    modelsToMigrate := []interface{}{
        &models.User{}, // Base models
    }
    modelsToMigrate = append(modelsToMigrate, registry.GetRegisteredModels()...)

    // Drop all tables
    for _, model := range modelsToMigrate {
        if err := db.Migrator().DropTable(model); err != nil {
            fmt.Printf("Failed to drop table: %v\n", err)
            return
        }
    }

    fmt.Println("Running fresh migrations...")
    if err := migrations.Migrate(db); err != nil {
        fmt.Printf("Failed to run migrations: %v\n", err)
        return
    }

    fmt.Println("Fresh migration completed successfully")
}

func runSpecificMigration(tableName string) {
    db := initDB()
    if db == nil {
        return
    }

    // Capitalize first letter for model name
    modelName := strings.Title(strings.ToLower(tableName))
    
    // Try to find the model in registry first
    var modelToMigrate interface{}
    if model, exists := registry.ModelRegistry[modelName]; exists {
        modelToMigrate = model
    } else {
        // Check if it's a base model (like User)
        switch modelName {
        case "User":
            modelToMigrate = &models.User{}
        default:
            fmt.Printf("Model %s not found\n", modelName)
            return
        }
    }

    fmt.Printf("Migrating %s table...\n", modelName)
    
    if err := db.AutoMigrate(modelToMigrate); err != nil {
        fmt.Printf("Failed to migrate %s: %v\n", modelName, err)
        return
    }

    fmt.Printf("Successfully migrated %s table\n", modelName)
}

func initDB() *gorm.DB {
    // Load configuration
    cfg, err := config.LoadConfig()
    if err != nil {
        fmt.Printf("Failed to load configuration: %v\n", err)
        return nil
    }

    // Initialize database
    db, err := config.InitDB(cfg)
    if err != nil {
        fmt.Printf("Failed to connect to database: %v\n", err)
        return nil
    }

    return db
}