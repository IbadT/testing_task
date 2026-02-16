package main

import (
	"fmt"
	"net/http"
	"testingtask/internal/config"
	"testingtask/internal/database"

	"testingtask/internal/delivery/http/middleware"
	v1 "testingtask/internal/delivery/http/v1"
	"testingtask/internal/repository"
	"testingtask/internal/service"
	"testingtask/internal/web/subscriptions"
	logger "testingtask/pkg"

	_ "testingtask/docs"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Subscription API
// @version 1.0
// @description API для управления подписками
// @host localhost:8081
// @BasePath /api
// @schemes http
func main() {
	e := echo.New()

	e.Use(middleware.RequestLoggerMiddleware)

	cfg, err := config.LoadConfig()
	if err != nil {
		panic("Failed to load config: " + err.Error())
	}

	appEnv := cfg.AppEnv
	logger.Init(appEnv == "development")

	db, err := database.InitDB(cfg)
	if err != nil {
		panic("Failed to init database: " + err.Error())
	}

	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	router := e.Group("/api")

	subRepo := repository.NewSubRepository(db)
	subService := service.NewSubService(subRepo)
	subHandler := v1.NewSubHandler(subService)

	subStrictHandler := subscriptions.NewStrictHandler(subHandler, nil)
	subscriptions.RegisterHandlers(router, subStrictHandler)

	port := fmt.Sprintf(":%s", cfg.PORT)

	if err := e.Start(port); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
