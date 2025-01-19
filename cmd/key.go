package cmd

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var keyCmd = &cobra.Command{
	Use:   "key:generate",
	Short: "Generate application key",
	Run: func(cmd *cobra.Command, args []string) {
		key := generateRandomKey(32)
		updateEnvFile("APP_KEY", key)
		fmt.Printf("Application key set successfully: %s\n", key)
	},
}

func init() {
	rootCmd.AddCommand(keyCmd)
}

func generateRandomKey(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(bytes)
}

func updateEnvFile(key, value string) {
	env, err := godotenv.Read(".env")
	if err != nil {
		// If .env doesn't exist, create a new one
		env = make(map[string]string)
	}

	env[key] = value

	content := ""
	for k, v := range env {
		content += fmt.Sprintf("%s=%s\n", k, v)
	}

	os.WriteFile(".env", []byte(content), 0644)
}
