package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"RESTAPI/internal/entity"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Connect устанавливает соединение с PostgreSQL
func Connect() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		"localhost", // Адрес базы данных
		"postgres",  // Пользователь
		"postgres",  // Пароль
		"distcomp",  // Название схемы/базы данных
		5432,        // Порт
	)

	// Подключение к базе данных
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Логирование успешного подключения
	log.Println("Successfully connected to PostgreSQL")

	// Автоматическое создание таблиц
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.WithContext(ctx).AutoMigrate(
		&entity.Writer{},
		&entity.News{},
		&entity.Message{},
		&entity.Mark{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate models: %w", err)
	}

	log.Println("Tables created successfully", entity.Message{}, entity.Writer{})

	return db, nil

}

func TestDatabaseConnection(db *gorm.DB) {
	var writer entity.Writer
	result := db.First(&writer)
	if result.Error != nil {
		log.Fatalf("Failed to fetch data: %v", result.Error)
	}
	log.Printf("Fetched writer: %+v", writer)
}
