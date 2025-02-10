package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/urlpick/urlpick-api/internal/app/url"
	"github.com/urlpick/urlpick-api/internal/pkg/config"
	"github.com/urlpick/urlpick-api/internal/pkg/db/mysql"
	"github.com/urlpick/urlpick-api/internal/pkg/middleware"
)

func main() {
	config.Load()

	mysqlDB, err := mysql.NewConnection()
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}
	defer mysqlDB.Close()

	urlRepo := mysql.NewURLRepository(mysqlDB)
	urlService := url.NewService(urlRepo)
	urlHandler := url.NewHandler(urlService)

	router := gin.Default()
	router.Use(middleware.CORS())
	router.Use(middleware.ErrorHandler())

	api := router.Group("/api/v1/")
	{
		urlHandler.RegisterRoutes(api)
	}

	if err := router.Run(":" + config.AppConfig.AppPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
