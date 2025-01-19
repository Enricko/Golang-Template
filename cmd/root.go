// cmd/root.go
package cmd

import (
    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "app",
    Short: "Your application CLI tool",
    Long:  `A command line interface for managing your Golang application`,
}

func Execute() error {  // Changed function signature to return error
    return rootCmd.Execute()  // Return the error from cobra's Execute
}