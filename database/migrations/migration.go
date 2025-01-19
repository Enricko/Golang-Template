// database/migrations/migrations.go
package migrations

import (
	"fmt"
	"golang-template/app/models"
	"golang-template/registry"
	"reflect"

	"gorm.io/gorm"
)

var ModelRegistry = make(map[string]interface{})

func RegisterModel(name string, model interface{}) {
    ModelRegistry[name] = model
}


func Migrate(db *gorm.DB) error {
    // Enable foreign key constraints for SQLite
    if db.Dialector.Name() == "sqlite" {
        err := db.Exec("PRAGMA foreign_keys = ON").Error
        if err != nil {
            return err
        }
    }

    // Initialize slice with base models
    modelsToMigrate := []interface{}{
        &models.User{}, // Base models
    }

    // Add all dynamically registered models
    modelsToMigrate = append(modelsToMigrate, registry.GetRegisteredModels()...)

    // Auto migrate all models
    for _, model := range modelsToMigrate {
        modelType := reflect.TypeOf(model).Elem()
        fmt.Printf("Migrating %s...\n", modelType.Name())
        
        if err := db.AutoMigrate(model); err != nil {
            return fmt.Errorf("failed to migrate %s: %v", modelType.Name(), err)
        }
    }

    return seedData(db)
}

func seedData(db *gorm.DB) error {
    var count int64
    db.Model(&models.User{}).Where("role = ?", "admin").Count(&count)
    
    if count == 0 {
        adminUser := models.User{
            Name:     "Admin",
            Email:    "admin@example.com",
            Password: "admin123",
            Role:     "admin",
        }
        
        if err := db.Create(&adminUser).Error; err != nil {
            return err
        }
    }
    
    return nil
}