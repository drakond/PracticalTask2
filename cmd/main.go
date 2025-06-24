package main

import (
	"log"
	"prtask2/internal/config"
)

func main() {
	// Загрузка конфигурации
	cfg := config.Load()

	// Запуск сервера
	log.Printf("🚀 %s starting on port %s", cfg.AppName, cfg.Port)
	app.Listen(cfg.Port)
}
