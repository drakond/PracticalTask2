package main

import (
	"context"
	"fmt"
	"log"
	"prtask2/internal/api"
	"prtask2/internal/config"
	"prtask2/internal/repo"
	"prtask2/internal/service"

	"github.com/jackc/pgx/v5"
)

func main() {
	// Загрузка конфигурации
	cfg := config.Load()

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer conn.Close(context.Background())

	// Инициализация зависимостей
	repository := repo.NewRepository(conn)
	taskService := service.NewService(repository)

	// Настройка роутов
	app := api.SetupRoutes(taskService)

	// Запуск сервера
	log.Printf("🚀 %s starting on port %s", cfg.AppName, cfg.Port)
	app.Listen(cfg.Port)
}
