package rss

import (
	"GoNews/pkg/models"
	"encoding/xml"
	"net/url"
	"time"

	"github.com/microcosm-cc/bluemonday"
)

// Структура тела новости, куда мы будем сохранять данные
type rssFeed struct {
	Channel struct {
		Items []struct {
			Title       string `xml:"title"`
			Link        string `xml:"link"`
			Description string `xml:"description"`
			Content     string `xml:"encoded"` //
			Content2    string `xml:"content"` // Возможные теги с содержимым тела новостей
			Summary     string `xml:"summary"` //
			PubDate     string `xml:"pubDate"`
		} `xml:"item"`
	} `xml:"channel"`
}

func ParseRSSFeed(b []byte) ([]models.News, error) {
	var feed rssFeed
	err := xml.Unmarshal(b, &feed)
	if err != nil {
		return nil, err
	}

	p := bluemonday.StripTagsPolicy() // Метод, который вызывается при каждой итерации для очистки контента от HTML тегов

	var news []models.News

	for _, item := range feed.Channel.Items {
		parsedTime, err := time.Parse(time.RFC1123Z, item.PubDate)
		text := item.Content // Многоэтапная проверка контента: если способ ничего не дал -> переключаемся на следующий
		if text == "" {
			text = item.Content2
		}
		if text == "" {
			text = item.Description
		}
		if text == "" {
			text = item.Summary
		}
		text = p.Sanitize(text) // Чистка от тегов с помощью готового решения blueMonday
		if err != nil {
			parsedTime = time.Now()
		}
		source, err := url.Parse(item.Link)
		if err != nil {
			source = &url.URL{} // Пустой url вместо nil в случае ошибки
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

// LoadAndParse загружает RSS по URL и парсит его в []models.News.
func LoadAndParse(url string, loader RSSLoader) ([]models.News, error) {
	raw, err := loader.LoadRSS(url)
	if err != nil {
		return nil, err
	}
	return ParseRSSFeed(raw)
}
