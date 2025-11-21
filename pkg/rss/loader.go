package rss

import (
	"io"
	"net/http"
)

type RSSLoader interface {
	LoadRSS(url string) ([]byte, error)
}
type HTTPRSSLoader struct{}

func (HTTPRSSLoader) LoadRSS(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
