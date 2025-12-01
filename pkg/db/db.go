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
	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Ошибка подключения к БД: %v", err)
		return nil, err
	}
	DB.Exec("TRUNCATE news RESTART IDENTITY CASCADE;") // реюзаем пространство, которое заняли в прошлом запуске
	err = DB.AutoMigrate(&models.News{})
	if err != nil {
		log.Printf("Ошибка Миграции БД:%v", err)
		return nil, err
	}
	log.Println("Успешное подключение к БД!")
	return DB, nil
}

// SaveNews сохраняет список новостей в БД, игнорируя дубликаты
func SaveNews(db *gorm.DB, news []models.News) error {
	for _, n := range news {
		err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&n).Error
		if err != nil {
			return err
		}
	}
	return nil
}
func GetLastNews(db *gorm.DB, limit int) ([]models.News, error) {
	var news []models.News
	query := `
        SELECT id, title, content, pub_time, link, source
        FROM (
            SELECT *,
                   ROW_NUMBER() OVER (
                       PARTITION BY source
                       ORDER BY pub_time DESC
                   ) AS rn
            FROM news
        ) t
        WHERE t.rn <= ?
        ORDER BY source, pub_time DESC;
    `
	err := db.Raw(query, limit).Scan(&news).Error
	if err != nil {
		return nil, err
	}
	return news, nil
}

func GetNewsByID(db *gorm.DB, id int) (models.News, error) {
	var news models.News
	err := db.First(&news, id).Error
	return news, err
}

func Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
