REST API для системы управления сущностями

**Описание**
Этот проект представляет собой RESTful API для управления сущностями: Writer, News, Mark и Message. Проект реализован на языке программирования Go с использованием фреймворка Echo. Он поддерживает основные CRUD-операции для каждой сущности и может быть легко расширен для добавления новых функций.

---

## Содержание
1. [Требования](#требования)
2. [Установка](#установка)
3. [Структура проекта](#структура-проекта)
4. [Использование](#использование)
5. [API Документация](#api-документация)
6. [Тестирование](#тестирование)
7. [Авторы](#авторы)

---

## Требования

Для запуска проекта вам понадобятся следующие зависимости:
- **Go**: версия 1.18 или выше
- **Echo Framework**: используется для создания сервера
- **Validator**: используется для валидации входных данных

Установите необходимые зависимости:
```bash
go get github.com/labstack/echo/v4
go get github.com/go-playground/validator/v10
```

---

## Установка

1. **Клонирование репозитория**:
   ```bash
   git clone https://github.com/username/your-repo.git
   cd your-repo
   ```

2. **Установка зависимостей**:
   ```bash
   go mod tidy
   ```

3. **Запуск приложения**:
   ```bash
   go run cmd/server/main.go
   ```

Приложение будет доступно по адресу: `http://localhost:8080`.

---

## Структура проекта

Проект организован по принципу слоев:

```
cmd/
└── server/
    └── main.go  # Точка входа в приложение
internal/
├── entity/      # Модели данных (Writer, News, Mark, Message)
│   ├── writer.go
│   ├── news.go
│   ├── mark.go
│   └── message.go
├── dto/         # DTO для обмена данными между клиентом и сервером
│   ├── writer_dto.go
│   ├── news_dto.go
│   ├── mark_dto.go
│   └── message_dto.go
├── storage/     # Хранилище данных (InMemory)
│   ├── repository.go
│   ├── writer_storage.go
│   ├── news_storage.go
│   ├── mark_storage.go
│   └── message_storage.go
├── service/     # Бизнес-логика
│   ├── writer_service.go
│   ├── news_service.go
│   ├── mark_service.go
│   └── message_service.go
├── handler/     # Обработчики HTTP-запросов
│   ├── writer_handler.go
│   ├── news_handler.go
│   ├── mark_handler.go
│   └── message_handler.go
└── validator/   # Валидация данных
    └── validator.go
```

---

## Использование

### Запуск сервера
Чтобы запустить сервер, выполните следующую команду:
```bash
go run cmd/server/main.go
```

По умолчанию сервер будет запущен на порту `24110` для проверки работоспособности на стороннем сервисе.

### Доступные маршруты
#### Writer
- **POST /api/v1.0/writers**: Создание писателя
- **GET /api/v1.0/writers/:id**: Получение писателя по ID
- **PUT /api/v1.0/writers/:id**: Обновление писателя
- **DELETE /api/v1.0/writers/:id**: Удаление писателя
- **GET /api/v1.0/writers**: Получение списка всех писателей

#### News
- **POST /api/v1.0/news**: Создание новости
- **GET /api/v1.0/news/:id**: Получение новости по ID
- **PUT /api/v1.0/news/:id**: Обновление новости
- **DELETE /api/v1.0/news/:id**: Удаление новости
- **GET /api/v1.0/news**: Получение списка всех новостей

#### Message
- **POST /api/v1.0/messages**: Создание сообщения
- **GET /api/v1.0/messages/:id**: Получение сообщения по ID
- **PUT /api/v1.0/messages/:id**: Обновление сообщения
- **DELETE /api/v1.0/messages/:id**: Удаление сообщения
- **GET /api/v1.0/messages**: Получение списка всех сообщений

#### Mark
- **POST /api/v1.0/marks**: Создание метки
- **GET /api/v1.0/marks/:id**: Получение метки по ID
- **PUT /api/v1.0/marks/:id**: Обновление метки
- **DELETE /api/v1.0/marks/:id**: Удаление метки
- **GET /api/v1.0/marks**: Получение списка всех меток

---

## API Документация

### Пример запроса к Writer
#### Создание писателя
```bash
curl -X POST http://localhost:8080/api/v1.0/writers \
    -H "Content-Type: application/json" \
    -d '{"login":"john_doe","password":"secure123","firstname":"John","lastname":"Doe"}'
```

#### Получение списка писателей
```bash
curl -X GET http://localhost:8080/api/v1.0/writers
```

---

## Тестирование

Для тестирования используйте инструменты, такие как Postman или curl. Также можно написать автоматические тесты с использованием Go Testing Framework.

Пример теста:
```go
package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestCreateWriter(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1.0/writers", strings.NewReader(`{"login":"test","password":"test123","firstname":"Test","lastname":"User"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := NewWriterHandler(nil) // Передайте сервис для тестирования
	if err := handler.Create(c); err != nil {
		t.Fatal(err)
	}

	if rec.Code != http.StatusCreated {
		t.Errorf("Expected status %v, got %v", http.StatusCreated, rec.Code)
	}
}
```

---

## Авторы

- **Danila Asepkov**: danilamail05@gmail.com  
  

---
