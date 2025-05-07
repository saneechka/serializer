# Serializer - универсальная библиотека для сериализации данных в Go

Библиотека предоставляет унифицированный интерфейс для сериализации/десериализации данных в различных форматах (JSON, TOML) и интеграцию с фреймворком Gin.

## Установка

```bash
go get github.com/saneechka/serializer
```

## Основные возможности

- Унифицированный интерфейс для работы с разными форматами
- Поддержка JSON и TOML форматов
- Интеграция с фреймворком Gin

## Использование

### Базовое использование

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/saneechka/serializer"
)

type User struct {
    ID    int    `json:"id" toml:"id"`
    Name  string `json:"name" toml:"name"`
    Email string `json:"email" toml:"email"`
}

func main() {
    // Создание сериализатора JSON
    jsonSerializer, err := serializer.New("json")
    if err != nil {
        log.Fatal(err)
    }
    
    user := User{
        ID:    1,
        Name:  "Иван",
        Email: "ivan@example.com",
    }
    
    // Сериализация в JSON
    jsonData, err := jsonSerializer.Marshal(user)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(string(jsonData))
    
    // Десериализация из JSON
    var decodedUser User
    err = jsonSerializer.Unmarshal(jsonData, &decodedUser)
    if err != nil {
        log.Fatal(err)
    }
    
    // Создание сериализатора TOML
    tomlSerializer, err := serializer.New("toml")
    if err != nil {
        log.Fatal(err)
    }
    
    // Сериализация в TOML
    tomlData, err := tomlSerializer.Marshal(user)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(string(tomlData))
}
```

### Интеграция с Gin

Библиотека предоставляет готовые функции для использования с фреймворком Gin:

```go
package main

import (
    "net/http"
    
    "github.com/gin-gonic/gin"
    mygin "github.com/saneechka/serializer/gin"
)

type User struct {
    ID       int    `json:"id" toml:"id"`
    Name     string `json:"name" toml:"name"`
    Email    string `json:"email" toml:"email"`
    Password string `json:"password" toml:"password"`
}

type UserService struct {
    // Методы для работы с пользователями...
}

func (s *UserService) Register(user *User) error {
    // Логика регистрации пользователя...
    return nil
}

type UserHandler struct {
    userService *UserService
}

// Пример обработчика Gin с использованием кастомной сериализации
func (h *UserHandler) Register(c *gin.Context) {
    var user User
    
    // Использование кастомной сериализации вместо встроенных методов Gin
    if err := mygin.MyBindJSON(c, &user); err != nil {
        mygin.MyJSON(c, http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.userService.Register(&user); err != nil {
        mygin.MyJSON(c, http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user.Password = "" // Не возвращаем пароль в ответе
    mygin.MyJSON(c, http.StatusCreated, gin.H{"message": "User registered successfully", "user": user})
}

// Пример обработчика Gin с использованием TOML
func (h *UserHandler) RegisterWithTOML(c *gin.Context) {
    var user User
    
    // Десериализация запроса в формате TOML
    if err := mygin.MyBindTOML(c, &user); err != nil {
        mygin.MyTOML(c, http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.userService.Register(&user); err != nil {
        mygin.MyTOML(c, http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user.Password = ""
    mygin.MyTOML(c, http.StatusCreated, gin.H{"message": "User registered successfully", "user": user})
}

func main() {
    router := gin.Default()
    
    userService := &UserService{}
    userHandler := &UserHandler{userService: userService}
    
    // Регистрация маршрутов
    router.POST("/register", userHandler.Register)
    router.POST("/register-toml", userHandler.RegisterWithTOML)
    
    router.Run(":8080")
}
```

## API библиотеки

### Интерфейс Serializer

```go
type Serializer interface {
    Marshal(v any) ([]byte, error)
    Unmarshal(data []byte, v any) error
    Format() string
}
```

### Основные функции

- `serializer.New(format string) (Serializer, error)` - создает новый сериализатор для указанного формата ("json" или "toml")
- `serializer.NewGin(format string) (*GinSerializer, error)` - создает новый сериализатор для использования с Gin

### Методы для Gin

- `gin.MyBindJSON(c *gin.Context, obj any) error` - десериализует JSON данные из запроса в объект
- `gin.MyBindTOML(c *gin.Context, obj any) error` - десериализует TOML данные из запроса в объект
- `gin.MyJSON(c *gin.Context, code int, obj any) error` - сериализует объект в JSON и отправляет ответ
- `gin.MyTOML(c *gin.Context, code int, obj any) error` - сериализует объект в TOML и отправляет ответ

## Обработка ошибок

Библиотека предоставляет предопределенные ошибки:
- `ErrUnsupportedFormat` - возвращается, если запрошенный формат сериализации не поддерживается

```go
if err == serializer.ErrUnsupportedFormat {
    log.Fatal("Формат не поддерживается")
}
```

## Пример применения в проекте

### Модификация существующего кода Gin

Если у вас уже есть обработчики Gin, вы можете легко заменить встроенные функции на кастомные:

Было:
```go
func (h *UserHandler) Register(c *gin.Context) {
    var user models.User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.userService.Register(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user.Password = ""
    c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "user": user})
}
```

Стало:
```go
func (h *UserHandler) Register(c *gin.Context) {
    var user models.User
    
    // Используем MyBindJSON вместо c.ShouldBindJSON
    if err := gin.MyBindJSON(c, &user); err != nil {
        // Используем MyJSON вместо c.JSON для ответа об ошибке
        gin.MyJSON(c, http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.userService.Register(&user); err != nil {
        // Используем MyJSON вместо c.JSON для ответа об ошибке
        gin.MyJSON(c, http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user.Password = ""
    // Используем MyJSON вместо c.JSON для успешного ответа
    gin.MyJSON(c, http.StatusCreated, gin.H{"message": "User registered successfully", "user": user})
}
```

## Тестирование

Библиотека включает набор тестов для всех реализованных сериализаторов. Для запуска тестов:

```bash
go test -v ./...
```