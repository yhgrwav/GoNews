# GoNews

RSS-агрегатор на Go. Собирает новости из RSS-лент и предоставляет REST API для доступа к ним.

## Установка зависимостей

```bash
go mod download
```

## Настройка

1. Убедитесь, что PostgreSQL запущен и доступен.

2. Отредактируйте файл `config.json`:

```json
{
  "db_host": "localhost",
  "db_port": 5432,
  "db_user": "postgres",
  "db_password": "ваш_пароль",
  "db_name": "postgres",
  "rss": [
    "https://vc.ru/rss/all",
    "https://www.technologyreview.com/feed/",
    "https://habr.com/ru/rss/best/daily/?fl=ru",
    "https://moex.com/export/news.aspx?cat=100",
    "https://techcrunch.com/feed/"
  ]
}
```

Укажите параметры подключения к базе данных и список RSS-источников.

## Запуск базы данных (Docker)

```bash
docker-compose up -d
```

## Запуск приложения

```bash
go run main.go
```

Сервер запустится на порту 8080.

## Использование API

### Получить последние новости

```bash
curl http://localhost:8080/news?limit=5
```

Параметр `limit` определяет количество новостей на каждый источник (по умолчанию 5).

### Принудительное обновление новостей

```bash
curl http://localhost:8080/news/refresh
```

Запускает обновление всех RSS-лент из конфигурации.

### Проверка состояния сервиса

```bash
curl http://localhost:8080/news/health
```

Возвращает статус подключения к базе данных.

## Сборка

```bash
go build -o gonews main.go
./gonews
```

## Остановка

Нажмите `Ctrl+C` для корректного завершения работы приложения.
