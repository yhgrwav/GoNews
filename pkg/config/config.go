package config

import (
	"GoNews/pkg/models"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// 1.LoadConfig() — отвечает только за чтение и разбор JSON-данных
// она не знает, откуда эти данные пришли — из файла, из интернета или из теста
// благодаря этому её легко проверять и использовать повторно
// 2.ReadConfig() — открывает файл config.json и передаёт его содержимое в LoadConfig()
// т.е. она просто загружает конфигурацию из файла при запуске программы
func LoadConfig(r io.Reader) (*models.Config, error) {
	var cfg models.Config
	if err := json.NewDecoder(r).Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func ReadConfig() (*models.Config, error) {
	data, err := os.Open("config.json")
	if err != nil {
		return nil, err
	}
	defer data.Close()
	cfg, err := LoadConfig(data)
	if err != nil {
		return nil, err
	}

	if envHost := os.Getenv("DB_HOST"); envHost != "" {
		cfg.DBHost = envHost
	}
	if envPort := os.Getenv("DB_PORT"); envPort != "" {
		var p int
		if _, err := fmt.Sscanf(envPort, "%d", &p); err == nil {
			cfg.DBPort = p
		}
	}
	if envUser := os.Getenv("DB_USER"); envUser != "" {
		cfg.DBUser = envUser
	}
	if envPassword := os.Getenv("DB_PASSWORD"); envPassword != "" {
		cfg.DBPassword = envPassword
	}
	if envName := os.Getenv("DB_NAME"); envName != "" {
		cfg.DBName = envName
	}
	return cfg, nil
}
