package worker

import (
	"GoNews/pkg/db"
	"GoNews/pkg/models"
	"GoNews/pkg/rss"
	"log"

	"gorm.io/gorm"
)

type Task struct {
	URL string
}

// Суть воркера - получить url из канала jobs, вызвать LoadAndParse, обработать ошибку и вернуть результат в канал results
func worker(id int, jobs <-chan Task, results chan<- []models.News, loader rss.RSSLoader) {
	for task := range jobs {
		news, err := rss.LoadAndParse(task.URL, loader)
		if err != nil {
			log.Printf("Worker %d: ошибка загрузки %s: %v", id, task.URL, err)
			results <- []models.News{}
			continue
		}
		results <- news
	}
}

// У нас есть 2 канала:
// jobs - канал, который по размерам соответствует кол-ву rss-лент в конфиге. Его функция в проекте - выдавать воркерам
// url'ы, которые воркеры обязаны передать как аргумент в вызове функции LoadAndParse, тем самым параллельно быстро
// получить свежие новости со всех источников не теряя в оптимизации
// results - канал, который принимает результат воркеров. он также равен по длине количеству источников в конфиге.
// роль results в проекте - безопасно передать результат дальше по заданному маршруту данных.

func StartWorkerPool(cfg *models.Config, DB *gorm.DB, loader rss.RSSLoader) {
	jobs := make(chan Task, len(cfg.RSS))
	results := make(chan []models.News, len(cfg.RSS))

	for i := 0; i < 5; i++ { // 5 горутин теоретически не перерасходуют ресурсы, соответственно этого хватит для абсолютного большинства сценариев пользования
		go worker(i, jobs, results, loader)
	}

	go func() {
		for _, url := range cfg.RSS {
			jobs <- Task{URL: url}
		}
		close(jobs)

		var allNews []models.News
		for i := 0; i < len(cfg.RSS); i++ {
			res := <-results
			allNews = append(allNews, res...)
		}

		err := db.SaveNews(DB, allNews)
		if err != nil {
			log.Printf("Ошибка сохранения данных: %v", err)
		} else {
			log.Println("Обновление RSS завершено")
		}
	}()
}
