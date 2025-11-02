package main

import (
	"GoNews/pkg/rss"
	"log"

	"GoNews/pkg/api"
	"GoNews/pkg/config" // ← правильный пакет
	"GoNews/pkg/db"
)

func main() {
	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatal(err)
	}

	database, err := db.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}
	news, err := rss.ParseRSS(cfg, rss.HTTPRSSLoader{})
	if err != nil {
		log.Fatal("Ошибка парсинга RSS:", err)
	}

	err = db.SaveNews(database, news)
	if err != nil {
		log.Fatal("Ошибка сохранения новостей:", err)
	}
	err = api.StartServer(database, cfg)
	if err != nil {
		log.Fatal(err)
	}

}
