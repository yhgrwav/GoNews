package api

import (
	"GoNews/pkg/db"
	"GoNews/pkg/models"
	"GoNews/pkg/rss"
	"GoNews/pkg/worker"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func StartServer(database *gorm.DB, cfg *models.Config, loader rss.RSSLoader) *http.Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// GET /news?limit=5
	router.GET("/news", func(c *gin.Context) {
		limitStr := c.DefaultQuery("limit", "5")

		n, err := strconv.Atoi(limitStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "limit must be an integer"})
			return
		}

		if n <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "limit must be > 0"})
			return
		}
		news, err := db.GetLastNews(database, n)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "failed to fetch news"})
			return
		}
		c.JSON(http.StatusOK, news)
	})

	// GET /news/refresh
	// вызывает горутину, отвечающую за явный вызов функции LoadAndParse
	router.GET("/news/refresh", func(c *gin.Context) {
		go worker.StartWorkerPool(c, cfg, database, loader)
		c.JSON(http.StatusOK, gin.H{"Status": "refresh started"})
	})
	// GET /news/health
	// Возвращает 2 результата: состояние текущего соединения и результат пинга БД
	router.GET("/news/health", func(c *gin.Context) {
		sqlDB, err := database.DB()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "database connection problem"})
			return
		}
		err = sqlDB.Ping()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "database is unreachable"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"Status": "OK"})
	})
	return &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
}
