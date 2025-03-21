// @title Song Library API
// @version 1.0
// @description API для управления библиотекой песен
// @host localhost:8080
// @BasePath /api/v1

package main

import (
	"github.com/Zorynix/song-library/internal/app"
)

const configPath = "config/config.yaml"

func main() {
	app.Run(configPath)
}
