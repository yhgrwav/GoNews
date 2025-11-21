package api

import (
	"GoNews/pkg/db"
	"GoNews/pkg/models"
	"GoNews/pkg/rss"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func StartServer(database *gorm.DB, cfg *models.Config, loader rss.RSSLoader) error {
	router := gin.Default()

	// GET /news?limit=10
	router.GET("/news", func(c *gin.Context) {
		limitStr := c.DefaultQuery("limit", "1")

		n, err := strconv.Atoi(limitStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "limit must be an integer"})
			return
		}

		if n <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "limit must be > 0"})
			return
		}
		news, err := db.GetLastNews(database, n)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch news"})
			return
		}
		c.JSON(http.StatusOK, news)
	})

	// GET /news/:id
	router.GET("/news/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id must be an integer"})
			return
		}

		news, err := db.GetNewsByID(database, id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "news not found"})
			return
		}
		c.JSON(http.StatusOK, news)
	})

	// GET /news/refresh
	router.GET("/news/refresh", func(c *gin.Context) {
		err := rss.UpdateAll(cfg.RSS, loader, database)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to refresh RSS"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	//health
	return router.Run(":8080")
}
