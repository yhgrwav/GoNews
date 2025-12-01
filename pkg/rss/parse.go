package rss

import (
	"GoNews/pkg/models"
	"encoding/xml"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/microcosm-cc/bluemonday"
)

// Структура тела новости для RSS-источников
type rssFeed struct {
	Channel struct {
		Items []struct {
			Title       string `xml:"title"`
			Link        string `xml:"link"`
			Description string `xml:"description"`
			PubDate     string `xml:"pubDate"`
		} `xml:"item"`
	} `xml:"channel"`
}

// clearContent выполняет очистку html-мусора
func clearContent(content string) string {
	if content == "" {
		return ""
	}

	// 1) var p удаляет html разметку с помощью готового решения bluemonday
	p := bluemonday.StrictPolicy()
	content = p.Sanitize(content)

	// 2) spaceCleaner убирает \n \t \r с помощью regexp
	spaceCleaner := regexp.MustCompile(`[\n\t\r]{2,}`)
	content = spaceCleaner.ReplaceAllString(content, " ")

	// 3) strings.Join соединяет слова одним пробелом
	content = strings.Join(strings.Fields(content), " ")

	// 4) в конце функция отдаёт контент с дополнительной очисткой начала и конца каждой строки
	return strings.TrimSpace(content)
}

// ParseRSSFeed преобразует XML RSS-ленты в []models.News
func ParseRSSFeed(b []byte) ([]models.News, error) {
	var feed rssFeed
	err := xml.Unmarshal(b, &feed)
	if err != nil {
		return nil, err
	}

	news := make([]models.News, 0, len(feed.Channel.Items))

	for _, item := range feed.Channel.Items {
		parsedTime, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			parsedTime = time.Now()
		}

		// Единственное корректное поле контента для всех выбранных источников
		text := clearContent(item.Description)

		// Определяем домен источника
		source, err := url.Parse(item.Link)
		if err != nil {
			source = &url.URL{}
		}
		n := models.News{
			Title:   item.Title,
			Content: text,
			PubTime: parsedTime,
			Link:    item.Link,
			Source:  source.Host,
		}
		news = append(news, n)
	}
	return news, nil
}

// LoadAndParse загружает RSS-ленту по URL и парсит её в []models.News.
func LoadAndParse(url string, loader RSSLoader) ([]models.News, error) {
	raw, err := loader.LoadRSS(url)
	if err != nil {
		return nil, err
	}
	return ParseRSSFeed(raw)
}
