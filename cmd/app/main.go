// @title Song Library API
// @version 1.0
// @description API для управления библиотекой песен
// @host localhost:8080
// @BasePath /api/v1

package main

import (
	"fmt"
	"os"

	"github.com/Zorynix/song-library/internal/app"
)

const configPath = "config/config.yaml"

func main() {
	fmt.Println("Starting application...")
	if err := app.Run(configPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
