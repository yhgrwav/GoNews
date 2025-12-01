package main

import (
	"GoNews/pkg/api"
	"GoNews/pkg/config"
	"GoNews/pkg/db"
	"GoNews/pkg/rss"
	"GoNews/pkg/worker"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/net/context"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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
	worker.StartWorkerPool(ctx, cfg, database, loader)

	// 5. Запускаем сервер в горутине
	srv := api.StartServer(database, cfg, loader)
	log.Println("Starting HTTP API on :8080")
	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			fmt.Println(err)
			return
		}
	}()
	if err != nil {
		log.Fatalf("failed to start api: %v", err)
	}
	// 6. Graceful shutdown
	shutdown := make(chan os.Signal)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	// main.go блокируется, ожидая сигнал от ОС
	<-shutdown
	// Вызываем cancel(), который закрывает канал ctx.Done() — все воркеры мгновенно выходят через select.
	cancel()
	log.Println("Обработчики новостей отключены")
	// Выключаем HTTP-сервер
	err = srv.Shutdown(ctx)
	if err != nil {
		log.Fatalf("Не удалось отключить HTTP-сервер: %v. HTTP-сервер будет отключен принудительно", err)
	} else {
		log.Println("HTTP-сервер отключен")
	}
	// Закрываем пул соединений с БД
	err = db.Close(database)
	if err != nil {
		log.Fatalf("Не удалось отключить БД: %v. БД будет отключена принудительно", err)
	} else {
		log.Println("Соединения с БД закрыто")
	}

}
