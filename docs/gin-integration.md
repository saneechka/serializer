# Интеграция с Gin

В этом документе описаны подходы к интеграции библиотеки `serializer` с фреймворком Gin.

## Основные функции для интеграции с Gin

Библиотека предоставляет следующие функции для работы с Gin:

- `MyBindJSON`: альтернатива стандартному `c.ShouldBindJSON`
- `MyBindTOML`: привязка данных из TOML формата
- `MyJSON`: альтернатива стандартному `c.JSON`
- `MyTOML`: ответ в формате TOML

## Быстрый старт

### Шаг 1: Импортируйте необходимые пакеты

```go
import (
    "github.com/gin-gonic/gin"
    mygin "github.com/saneechka/serializer/gin"
)
```

### Шаг 2: Используйте кастомные функции вместо встроенных

```go
func YourHandler(c *gin.Context) {
    var data YourStruct
    
    // Чтение данных запроса
    if err := mygin.MyBindJSON(c, &data); err != nil {
        mygin.MyJSON(c, http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    // Обработка данных...
    
    // Отправка ответа
    mygin.MyJSON(c, http.StatusOK, gin.H{"data": data})
}
```

## Полная замена обработчиков в существующем проекте

Если у вас большой проект с множеством обработчиков, вы можете заменить все обращения к встроенным методам Gin на ваши кастомные методы. Для этого рекомендуется использовать следующий подход:

1. Определите общие обработчики ошибок:

```go
func ErrorResponse(c *gin.Context, status int, err error) {
    mygin.MyJSON(c, status, gin.H{"error": err.Error()})
}

func SuccessResponse(c *gin.Context, status int, data interface{}) {
    mygin.MyJSON(c, status, data)
}
```

2. Используйте эти функции во всех обработчиках:

```go
func (h *UserHandler) Register(c *gin.Context) {
    var user models.User
    
    if err := mygin.MyBindJSON(c, &user); err != nil {
        ErrorResponse(c, http.StatusBadRequest, err)
        return
    }

    if err := h.userService.Register(&user); err != nil {
        ErrorResponse(c, http.StatusBadRequest, err)
        return
    }

    user.Password = ""
    SuccessResponse(c, http.StatusCreated, gin.H{
        "message": "User registered successfully", 
        "user": user,
    })
}
```

## Поддержка нескольких форматов в одном обработчике

Вы можете реализовать обработчики, которые поддерживают несколько форматов сериализации в зависимости от заголовка `Content-Type`:

```go
func MultiFormatHandler(c *gin.Context) {
    var data YourStruct
    contentType := c.GetHeader("Content-Type")
    
    var err error
    switch contentType {
    case "application/json":
        err = mygin.MyBindJSON(c, &data)
    case "application/toml":
        err = mygin.MyBindTOML(c, &data)
    default:
        mygin.MyJSON(c, http.StatusUnsupportedMediaType, gin.H{
            "error": "Unsupported media type",
        })
        return
    }
    
    if err != nil {
        mygin.MyJSON(c, http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    // Обработка данных...
    
    // Ответ в том же формате, что и запрос
    switch contentType {
    case "application/json":
        mygin.MyJSON(c, http.StatusOK, gin.H{"data": data})
    case "application/toml":
        mygin.MyTOML(c, http.StatusOK, gin.H{"data": data})
    }
}
```

## Мидлвары для автоматического определения формата

Вы можете создать мидлвару, которая будет автоматически определять формат данных и устанавливать соответствующий обработчик:

```go
func FormatDetector() gin.HandlerFunc {
    return func(c *gin.Context) {
        contentType := c.GetHeader("Content-Type")
        
        // Устанавливаем формат в контексте
        switch contentType {
        case "application/json":
            c.Set("format", "json")
        case "application/toml":
            c.Set("format", "toml")
        default:
            c.Set("format", "json") // По умолчанию используем JSON
        }
        
        c.Next()
    }
}

// Использование:
r := gin.Default()
r.Use(FormatDetector())
```

Затем в обработчиках:

```go
func Handler(c *gin.Context) {
    format, _ := c.Get("format")
    
    var data YourStruct
    if format == "json" {
        mygin.MyBindJSON(c, &data)
    } else if format == "toml" {
        mygin.MyBindTOML(c, &data)
    }
    
    // ...
}
```

## Советы по оптимизации

### Создание сериализаторов один раз

Для повышения производительности рекомендуется создавать сериализаторы один раз при запуске приложения:

```go
type App struct {
    jsonSerializer serializer.Serializer
    tomlSerializer serializer.Serializer
}

func NewApp() *App {
    jsonSerializer, _ := serializer.New("json")
    tomlSerializer, _ := serializer.New("toml")
    
    return &App{
        jsonSerializer: jsonSerializer,
        tomlSerializer: tomlSerializer,
    }
}

func (a *App) SetupRoutes() *gin.Engine {
    r := gin.Default()
    
    r.POST("/data", a.handleData)
    
    return r
}

func (a *App) handleData(c *gin.Context) {
    // Использование предварительно созданных сериализаторов
    // ...
}
```

### Обработка ошибок

Не забывайте всегда обрабатывать ошибки, возвращаемые функциями сериализации:

```go
if err := mygin.MyJSON(c, http.StatusOK, data); err != nil {
    // Обработка ошибки сериализации
    c.String(http.StatusInternalServerError, "Error serializing response")
    return
}
```

## Примеры тестирования обработчиков

```go
func TestUserHandler_Register(t *testing.T) {
    router := gin.Default()
    
    userService := &UserService{}
    userHandler := &UserHandler{userService: userService}
    
    router.POST("/register", userHandler.Register)
    
    // Создаем тестовые данные
    user := User{
        Name:     "Test User",
        Email:    "test@example.com",
        Password: "password123",
    }
    
    jsonData, _ := json.Marshal(user)
    req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonData))
    req.Header.Set("Content-Type", "application/json")
    
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    // Проверяем ответ
    if w.Code != http.StatusCreated {
        t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
    }
    
    var response map[string]interface{}
    if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
        t.Fatal(err)
    }
    
    if msg, ok := response["message"]; !ok || msg != "User registered successfully" {
        t.Errorf("Expected success message, got %v", response)
    }
}
```