package config

import (
	"GoNews/pkg/models"
	"encoding/json"
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
	return LoadConfig(data)
}
func Validate(c *models.Config) error {
	if c.RequestPeriod <= 0 {
		return fmt.Errorf("RequestPeriod must be > 0")
	}
	if len(c.RSS) == 0 {
		return fmt.Errorf("no RSS feeds specified")
	}
	if len(c.DBHost) == 0 {
		return fmt.Errorf("no DBHost specified")
	}
	if len(c.DBUser) == 0 {
		return fmt.Errorf("no DBUser specified")
	}
	if len(c.DBPassword) == 0 {
		return fmt.Errorf("no DBPassword specified")
	}
	if len(c.RSS) == 0 {
		return fmt.Errorf("no RSS feeds specified")
	}
	return nil
}
