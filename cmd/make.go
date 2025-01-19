// cmd/make.go
package cmd

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"
    "text/template"

    "github.com/spf13/cobra"
)

// Flag variables
var (
    makeModelFlag bool
    forceFlag     bool
)

var makeModelCmd = &cobra.Command{
    Use:   "make:model [name]",
    Short: "Create a new model",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        createModel(args[0])
    },
}

var makeControllerCmd = &cobra.Command{
    Use:   "make:controller [name]",
    Short: "Create a new controller",
    Long:  `Create a new controller, optionally with a matching model using the -m flag`,
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        if makeModelFlag {
            // If -m flag is provided, create model first
            createModel(args[0])
        }
        createController(args[0])
    },
}

func init() {
    // Add flags to commands
    makeControllerCmd.Flags().BoolVarP(&makeModelFlag, "model", "m", false, "Create a model along with the controller")
    makeControllerCmd.Flags().BoolVarP(&forceFlag, "force", "f", false, "Force create/overwrite files")
    makeModelCmd.Flags().BoolVarP(&forceFlag, "force", "f", false, "Force create/overwrite files")
    
    // Add commands to root command
    rootCmd.AddCommand(makeModelCmd)
    rootCmd.AddCommand(makeControllerCmd)
}

// Templates
const modelTemplate = `package models

import (
    "gorm.io/gorm"
    "golang-template/registry"
)

type {{.Name}} struct {
    gorm.Model
    // Add your fields here
}

func init() {
    registry.RegisterModel("{{.Name}}", &{{.Name}}{})
}

// Implement model methods here
func (m *{{.Name}}) BeforeCreate(tx *gorm.DB) error {
    // Add any pre-create logic here
    return nil
}
`

const controllerTemplate = `package controllers

import (
    "net/http"
    "{{.ModuleName}}/app/models"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

type {{.Name}}Controller struct {
    db *gorm.DB
}

func New{{.Name}}Controller(db *gorm.DB) *{{.Name}}Controller {
    return &{{.Name}}Controller{db: db}
}

// Index retrieves all {{.PluralName}}
func (c *{{.Name}}Controller) Index(ctx *gin.Context) {
    var {{.LowerPluralName}} []models.{{.Name}}
    result := c.db.Find(&{{.LowerPluralName}})
    if result.Error != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }
    ctx.JSON(http.StatusOK, {{.LowerPluralName}})
}

// Show retrieves a single {{.Name}}
func (c *{{.Name}}Controller) Show(ctx *gin.Context) {
    id := ctx.Param("id")
    var {{.LowerName}} models.{{.Name}}
    
    if err := c.db.First(&{{.LowerName}}, id).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            ctx.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
            return
        }
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    ctx.JSON(http.StatusOK, {{.LowerName}})
}

// Store creates a new {{.Name}}
func (c *{{.Name}}Controller) Store(ctx *gin.Context) {
    var {{.LowerName}} models.{{.Name}}
    if err := ctx.ShouldBindJSON(&{{.LowerName}}); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    if err := c.db.Create(&{{.LowerName}}).Error; err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    ctx.JSON(http.StatusCreated, {{.LowerName}})
}

// Update modifies an existing {{.Name}}
func (c *{{.Name}}Controller) Update(ctx *gin.Context) {
    id := ctx.Param("id")
    var {{.LowerName}} models.{{.Name}}
    
    if err := c.db.First(&{{.LowerName}}, id).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            ctx.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
            return
        }
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    if err := ctx.ShouldBindJSON(&{{.LowerName}}); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    if err := c.db.Save(&{{.LowerName}}).Error; err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    ctx.JSON(http.StatusOK, {{.LowerName}})
}

// Delete removes a {{.Name}}
func (c *{{.Name}}Controller) Delete(ctx *gin.Context) {
    id := ctx.Param("id")
    var {{.LowerName}} models.{{.Name}}
    
    if err := c.db.First(&{{.LowerName}}, id).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            ctx.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
            return
        }
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    if err := c.db.Delete(&{{.LowerName}}).Error; err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    ctx.JSON(http.StatusOK, gin.H{"message": "Record deleted successfully"})
}
`

type TemplateData struct {
    Name            string // Original name
    LowerName       string // Lowercase name
    PluralName      string // Pluralized name
    LowerPluralName string // Lowercase pluralized name
    ModuleName      string // Module name from go.mod
}

// Update createModel function to also create a migration file if needed
func createModel(name string) {
    data := prepareTemplateData(name)
    
    // Check if file exists and force flag is not set
    if !forceFlag && fileExists("app/models/" + strings.ToLower(name) + ".go") {
        fmt.Printf("Model file already exists. Use -f flag to overwrite.\n")
        return
    }
    
    err := generateFile("app/models", strings.ToLower(name)+".go", modelTemplate, data)
    if err != nil {
        fmt.Printf("Error creating model: %v\n", err)
        return
    }

    fmt.Printf("Model created successfully: app/models/%s.go\n", strings.ToLower(name))
    fmt.Printf("Model will be automatically migrated on next server start\n")
}

func createController(name string) {
    data := prepareTemplateData(name)
    
    // Check if file exists and force flag is not set
    if !forceFlag && fileExists("app/controllers/" + strings.ToLower(name) + "_controller.go") {
        fmt.Printf("Controller file already exists. Use -f flag to overwrite.\n")
        return
    }
    
    err := generateFile("app/controllers", strings.ToLower(name)+"_controller.go", controllerTemplate, data)
    if err != nil {
        fmt.Printf("Error creating controller: %v\n", err)
        return
    }

    fmt.Printf("Controller created successfully: app/controllers/%s_controller.go\n", strings.ToLower(name))
}

func prepareTemplateData(name string) TemplateData {
    return TemplateData{
        Name:            strings.Title(name),
        LowerName:       strings.ToLower(name),
        PluralName:      strings.Title(name) + "s", // Simple pluralization
        LowerPluralName: strings.ToLower(name) + "s",
        ModuleName:      getModuleName(),
    }
}

func generateFile(dir, filename, tmpl string, data TemplateData) error {
    // Ensure directory exists
    if err := os.MkdirAll(dir, 0755); err != nil {
        return fmt.Errorf("failed to create directory: %v", err)
    }

    // Parse template
    t, err := template.New("template").Parse(tmpl)
    if err != nil {
        return fmt.Errorf("failed to parse template: %v", err)
    }

    // Create file
    filepath := filepath.Join(dir, filename)
    file, err := os.Create(filepath)
    if err != nil {
        return fmt.Errorf("failed to create file: %v", err)
    }
    defer file.Close()

    // Execute template
    if err := t.Execute(file, data); err != nil {
        return fmt.Errorf("failed to execute template: %v", err)
    }

    return nil
}

func getModuleName() string {
    // Read go.mod file to get module name
    data, err := os.ReadFile("go.mod")
    if err != nil {
        return "your-module-name"
    }

    lines := strings.Split(string(data), "\n")
    for _, line := range lines {
        if strings.HasPrefix(line, "module ") {
            return strings.TrimSpace(strings.TrimPrefix(line, "module "))
        }
    }

    return "your-module-name"
}

func fileExists(path string) bool {
    _, err := os.Stat(path)
    return !os.IsNotExist(err)
}