package config

import (
	"strings"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	testjson := `{
	"db_host": "localhost",
  "db_port": 5432,
  "db_user": "postgres",
  "db_password": "1",
  "db_name": "postgres",
  "rss": ["https://lenta.ru/rss"]
}`
	cfg, err := LoadConfig(strings.NewReader(testjson))
	if err != nil {
		t.Fatalf("LoadConfig вернул ошибку: %v", err)
	}
	if cfg.DBHost != "localhost" {
		t.Fatalf("Ожидался localhost, получен %s", cfg.DBHost)
	}
	if cfg.DBPort != 5432 {
		t.Fatalf("Ожидался порт 5432, получен %d", cfg.DBPort)
	}
	if len(cfg.RSS) < 1 {
		t.Fatalf("Ошибка получения rss ссылки: %s", cfg.RSS)
	}
}
