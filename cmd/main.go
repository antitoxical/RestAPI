package main

import (
	"RESTAPI/internal/handler"
	"RESTAPI/internal/service"
	"RESTAPI/internal/storage"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	// Инициализация хранилищ
	writerStorage := storage.NewWriterStorage()
	newsStorage := storage.NewNewsStorage()
	messageStorage := storage.NewMessageStorage()
	markStorage := storage.NewMarkStorage()

	// Инициализация сервисов
	writerService := service.NewWriterService(writerStorage)
	newsService := service.NewNewsService(newsStorage)
	messageService := service.NewMessageService(messageStorage)
	markService := service.NewMarkService(markStorage)

	// Инициализация обработчиков
	writerHandler := handler.NewWriterHandler(writerService)
	newsHandler := handler.NewNewsHandler(newsService)
	messageHandler := handler.NewMessageHandler(messageService)
	markHandler := handler.NewMarkHandler(markService)

	// Маршруты для Writer
	e.POST("/api/v1.0/writers", writerHandler.Create)
	e.GET("/api/v1.0/writers/:id", writerHandler.GetById)
	e.PUT("/api/v1.0/writers", writerHandler.Update)
	e.DELETE("/api/v1.0/writers/:id", writerHandler.Delete)
	e.GET("/api/v1.0/writers", writerHandler.GetAll)

	// Маршруты для News
	e.POST("/api/v1.0/news", newsHandler.Create)
	e.GET("/api/v1.0/news/:id", newsHandler.GetById)
	e.PUT("/api/v1.0/news", newsHandler.Update)
	e.DELETE("/api/v1.0/news/:id", newsHandler.Delete)
	e.GET("/api/v1.0/news", newsHandler.GetAll)

	// Маршруты для Message
	e.POST("/api/v1.0/messages", messageHandler.Create)
	e.GET("/api/v1.0/messages/:id", messageHandler.GetById)
	e.PUT("/api/v1.0/messages", messageHandler.Update)
	e.DELETE("/api/v1.0/messages/:id", messageHandler.Delete)
	e.GET("/api/v1.0/messages", messageHandler.GetAll)

	// Маршруты для Mark
	e.POST("/api/v1.0/marks", markHandler.Create)
	e.GET("/api/v1.0/marks/:id", markHandler.GetById)
	e.PUT("/api/v1.0/marks", markHandler.Update)
	e.DELETE("/api/v1.0/marks/:id", markHandler.Delete)
	e.GET("/api/v1.0/marks", markHandler.GetAll)

	e.Logger.Fatal(e.Start(":24110"))
}
