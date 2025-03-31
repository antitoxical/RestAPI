package main

import (
	"RESTAPI/db"
	"RESTAPI/internal/entity"
	"RESTAPI/internal/handler"
	"RESTAPI/internal/repository"
	"RESTAPI/internal/service"
	"log"

	"github.com/labstack/echo/v4"
)

func main() {

	db, err := db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	e := echo.New()

	// Инициализация хранилищ
	writerRepo := repository.NewWriterRepository(db)
	log.Println(entity.Writer{})
	newsRepo := repository.NewNewsRepository(db)
	markRepo := repository.NewMarkRepository(db)
	messageRepo := repository.NewMessageRepository(db)

	// Создание сервисов
	writerService := service.NewWriterService(writerRepo)
	newsService := service.NewNewsService(newsRepo)
	markService := service.NewMarkService(markRepo)
	messageService := service.NewMessageService(messageRepo)

	// Создание обработчиков
	writerHandler := handler.NewWriterHandler(writerService)
	newsHandler := handler.NewNewsHandler(newsService)
	markHandler := handler.NewMarkHandler(markService)
	messageHandler := handler.NewMessageHandler(messageService)

	log.Println("Tables created successfully", entity.Message{}, entity.Writer{})

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
