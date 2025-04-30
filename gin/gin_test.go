package gin

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

type TestUser struct {
	ID    int    `json:"id" toml:"id"`
	Name  string `json:"name" toml:"name"`
	Email string `json:"email" toml:"email"`
}

func TestBindJSON(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	testUser := TestUser{ID: 1, Name: "Иван", Email: "ivan@example.com"}
	jsonData, _ := json.Marshal(testUser)

	req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	var result TestUser
	err := MyBindJSON(c, &result)
	if err != nil {
		t.Fatalf("BindJSON() error = %v", err)
	}

	if result.ID != testUser.ID || result.Name != testUser.Name || result.Email != testUser.Email {
		t.Errorf("BindJSON() = %v, want %v", result, testUser)
	}
}

func TestBindTOML(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	tomlData := []byte(`
id = 1
name = "Иван"
email = "ivan@example.com"
`)

	req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBuffer(tomlData))
	req.Header.Set("Content-Type", "application/toml")
	c.Request = req

	var result TestUser
	err := MyBindTOML(c, &result)
	if err != nil {
		t.Fatalf("BindTOML() error = %v", err)
	}

	expected := TestUser{ID: 1, Name: "Иван", Email: "ivan@example.com"}
	if result.ID != expected.ID || result.Name != expected.Name || result.Email != expected.Email {
		t.Errorf("BindTOML() = %v, want %v", result, expected)
	}
}

func TestTOML(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	testUser := TestUser{ID: 1, Name: "Иван", Email: "ivan@example.com"}

	err := MyTOML(c, http.StatusOK, testUser)
	if err != nil {
		t.Fatalf("TOML() error = %v", err)
	}

	if w.Code != http.StatusOK {
		t.Errorf("TOML() status = %v, want %v", w.Code, http.StatusOK)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/toml" {
		t.Errorf("TOML() Content-Type = %v, want %v", contentType, "application/toml")
	}

	var result TestUser
	c.Request = httptest.NewRequest(http.MethodPost, "/test", w.Body)
	err = MyBindTOML(c, &result)
	if err != nil {
		t.Fatalf("TOML() invalid response body: %v", err)
	}

	if result.ID != testUser.ID || result.Name != testUser.Name || result.Email != testUser.Email {
		t.Errorf("TOML() = %v, want %v", result, testUser)
	}
}
