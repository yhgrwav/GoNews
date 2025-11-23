package main

import (
	"GoNews/pkg/api"
	"GoNews/pkg/config"
	"GoNews/pkg/db"
	"GoNews/pkg/rss"
	"GoNews/pkg/worker"
	"log"
)

func main() {

	// 1. Загружаем конфиг
	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	// 2. Подключаемся к базе
	database, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// 3. Создаём загрузчик RSS
	loader := rss.HTTPRSSLoader{}

	// 4. Запускаем worker pool для фонового обновления RSS
	worker.StartWorkerPool(cfg, database, loader)

	// 5. Запускаем сервер
	log.Println("Starting HTTP API on :8080")
	err = api.StartServer(database, cfg, loader)
	if err != nil {
		log.Fatalf("failed to start api: %v", err)
	}
}
