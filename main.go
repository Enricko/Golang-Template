package main

import (
    "log"
    "golang-template/cmd"
)

func main() {
    if err := cmd.Execute(); err != nil {
        log.Fatal("Error executing command:", err)
    }
}