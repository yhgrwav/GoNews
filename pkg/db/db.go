package db

import (
	"GoNews/pkg/models"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// connect + automigrate
func Connect(cfg *models.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=UTC",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Ошибка подключения к БД: %v", err)
		return nil, err
	}
	err = db.AutoMigrate(&models.News{})
	if err != nil {
		log.Printf("Ошибка Миграции БД:%v", err)
		return nil, err
	}
	log.Println("Успешное подключение к БД!")
	return db, nil
}
func SaveNews(db *gorm.DB, news []models.News) error {
	for _, n := range news {
		err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&n).Error
		if err != nil {
			return err
		}
	}
	return nil
}
func GetLastNews(db *gorm.DB, n int) ([]models.News, error) {
	var news []models.News
	err := db.Order("pub_time DESC").Limit(n).Find(&news).Error
	if err != nil {
		return nil, fmt.Errorf("Ошибка при получении новостей: %v", err)
	}
	return news, nil
}
