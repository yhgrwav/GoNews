package api

import (
	"GoNews/pkg/db"
	"GoNews/pkg/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func StartServer(database *gorm.DB, cfg *models.Config) error {
	router := gin.Default()
	router.GET("/news/:n", func(c *gin.Context) {
		n, err := strconv.Atoi(c.Param("n"))
		if err != nil || n <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "некорректный параметр"})
			return
		}
		news, err := db.GetLastNews(database, n)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Ошибка": "не удалось получить новости"})
		}
		c.JSON(http.StatusOK, news)
	})

	err := router.Run(":8080")
	if err != nil {
		return err
	}
	return nil
}
