package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/urlpick/urlpick-api/internal/app/url"
	"github.com/urlpick/urlpick-api/internal/pkg/config"
	"github.com/urlpick/urlpick-api/internal/pkg/db/mysql"
	"github.com/urlpick/urlpick-api/internal/pkg/middleware"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	config.Load()

	mysqlDB, err := mysql.NewConnection()
	if err != nil {
		slog.Error("failed to connect to MySQL", "error", err)
		os.Exit(1)
	}
	defer mysqlDB.Close()

	urlRepo := mysql.NewURLRepository(mysqlDB)
	urlService := url.NewService(urlRepo)
	urlHandler := url.NewHandler(urlService)

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.Logger())
	router.SetTrustedProxies(nil)
	router.Use(middleware.CORS())
	router.Use(middleware.ErrorHandler())

	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	router.GET("/readyz", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()
		if err := mysqlDB.PingContext(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unavailable"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := router.Group("/api/v1/")
	{
		urlHandler.RegisterRoutes(api)
	}

	srv := &http.Server{
		Addr:              ":" + config.AppConfig.AppPort,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("server shutdown error", "error", err)
	}
}
