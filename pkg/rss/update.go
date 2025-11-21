package rss

import (
	"GoNews/pkg/db"

	"gorm.io/gorm"
)

// UpdateAll обходит все RSS-ленты из конфигурации и сохраняет новости в БД
func UpdateAll(rss []string, loader RSSLoader, database *gorm.DB) error {
	for _, link := range rss {
		news, err := LoadAndParse(link, loader)
		if err != nil {
			return err
		}
		err = db.SaveNews(database, news)
		if err != nil {
			return err
		}
	}
	return nil
}
