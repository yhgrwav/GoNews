package db

import (
	"GoNews/pkg/models"
	"testing"
)

func TestConnect(t *testing.T) {
	cfg := &models.Config{
		DBHost:     "localhost",
		DBPort:     5432,
		DBUser:     "postgres",
		DBPassword: "1",
		DBName:     "postgres",
	}
	result, err := Connect(cfg)
	if err != nil {
		t.Error("Ошибка подключения: ", err)
	}
	if result == nil {
		t.Fatal("Ожидался указатель *gorm.DB, получен nil")
	}

}
