package models

import "time"

type Config struct {
	RSS           []string `json:"rss"`
	RequestPeriod int      `json:"request_period"`
	NewsAmount    int      `json:"NewsAmount"`
	DBHost        string   `json:"db_host"`
	DBPort        int      `json:"db_port"`
	DBUser        string   `json:"db_user"`
	DBPassword    string   `json:"db_password"`
	DBName        string   `json:"db_name"`
	APIPort       string   `json:"api_port"`
}
type News struct {
	ID        uint      `gorm:"primaryKey"`
	Title     string    `gorm:"type:text;not null"`
	Content   string    `gorm:"type:text"`
	PubTime   time.Time `gorm:"index"`
	Link      string    `gorm:"uniqueIndex;type:text"`
	Source    string    `gorm:"type:text"`
	CreatedAt time.Time `gorm:"index"`
}
type RSS struct {
	Channel struct {
		Items []struct {
			Title       string `xml:"title"`
			Link        string `xml:"link"`
			Description string `xml:"description"`
			PubDate     string `xml:"pubDate"`
		} `xml:"item"`
	} `xml:"channel"`
}
