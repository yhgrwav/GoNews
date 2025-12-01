package worker

import (
	"GoNews/pkg/db"
	"GoNews/pkg/models"
	"GoNews/pkg/rss"
	"log"

	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type Task struct {
	URL string
}

// Суть воркера - получить url из канала jobs, вызвать LoadAndParse,
// обработать ошибку и вернуть результат в канал results
func worker(ctx context.Context, id int, jobs <-chan Task, results chan<- []models.News, loader rss.RSSLoader) {
	for {
		select {
		// если из main вызывается graceful shutdown — выходим
		case <-ctx.Done():
			log.Printf("Worker %d: остановлен", id)
			return

		// основная логика: получаем URL, грузим RSS, отдаём результат
		case task, ok := <-jobs:
			if !ok {
				return
			}

			news, err := rss.LoadAndParse(task.URL, loader)
			if err != nil {
				log.Printf("Worker %d: ошибка загрузки %s: %v", id, task.URL, err)
				results <- []models.News{}
				continue
			}
			results <- news
		}
	}
}

// Запускает обновление RSS
func StartWorkerPool(ctx context.Context, cfg *models.Config, DB *gorm.DB, loader rss.RSSLoader) {

	// У нас есть 2 канала:
	//
	// jobs - канал, который по размерам соответствует кол-ву rss-лент в конфиге. Его функция в проекте - выдавать воркерам
	// url'ы, которые воркеры обязаны передать как аргумент в вызове функции LoadAndParse, тем самым параллельно быстро
	// получить свежие новости со всех источников не теряя в оптимизации
	jobs := make(chan Task, len(cfg.RSS))
	// results - канал, который принимает результат воркеров. он также равен по длине количеству источников в конфиге.
	// Роль results в проекте - безопасно передать результат дальше по заданному маршруту данных.
	results := make(chan []models.News, len(cfg.RSS))

	// запускаем пул воркеров
	for i := 0; i < 5; i++ { // 5 горутин достаточно для любых сценариев
		go worker(ctx, i, jobs, results, loader)
	}

	// Отправление задач
	for _, url := range cfg.RSS {
		select {
		case <-ctx.Done(): // graceful shutdown во время refresh
			log.Println("Обновление RSS отменено во время shutdown")
			return
		case jobs <- Task{URL: url}:
		}
	}
	close(jobs)

	// Сбор результатов
	var allNews []models.News
	for i := 0; i < len(cfg.RSS); i++ {
		select {
		case <-ctx.Done(): // graceful shutdown
			log.Println("Сбор результатов отменён во время shutdown")
			return
		case res := <-results:
			allNews = append(allNews, res...)
		}
	}

	// Сохранение в БД
	err := db.SaveNews(DB, allNews)
	if err != nil {
		log.Printf("Ошибка сохранения данных: %v", err)
	} else {
		log.Println("Обновление RSS завершено")
	}
}
