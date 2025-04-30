package toml

import (
	"reflect"
	"testing"
)

func TestTOMLSerializer(t *testing.T) {
	type TestStruct struct {
		String  string   `toml:"string"`
		Integer int      `toml:"integer"`
		Float   float64  `toml:"float"`
		Boolean bool     `toml:"boolean"`
		Array   []string `toml:"array"`
		Nested  struct {
			Field string `toml:"field"`
		} `toml:"nested"`
	}

	serializer := New()

	if format := serializer.Format(); format != "TOML" {
		t.Errorf("Format() = %v, want %v", format, "TOML")
	}

	original := TestStruct{
		String:  "тест",
		Integer: 42,
		Float:   3.14,
		Boolean: true,
		Array:   []string{"один", "два", "три"},
		Nested: struct {
			Field string `toml:"field"`
		}{
			Field: "вложенное поле",
		},
	}

	data, err := serializer.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	var result TestStruct
	err = serializer.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	if !reflect.DeepEqual(original, result) {
		t.Errorf("Unmarshal() = %v, want %v", result, original)
	}

	err = serializer.Unmarshal([]byte("invalid = toml]"), &result)
	if err == nil {
		t.Error("Unmarshal() с неверным TOML должен возвращать ошибку")
	}
}

func TestTOMLComplexStructures(t *testing.T) {
	type Config struct {
		Server struct {
			Host string `toml:"host"`
			Port int    `toml:"port"`
			TLS  struct {
				Enabled     bool   `toml:"enabled"`
				Certificate string `toml:"certificate"`
				Key         string `toml:"key"`
			} `toml:"tls"`
		} `toml:"server"`
		Database struct {
			Driver   string            `toml:"driver"`
			Host     string            `toml:"host"`
			Username string            `toml:"username"`
			Password string            `toml:"password"`
			Options  map[string]string `toml:"options"`
		} `toml:"database"`
		Logging struct {
			Level  string   `toml:"level"`
			Output []string `toml:"output"`
		} `toml:"logging"`
	}

	serializer := New()

	original := Config{}
	original.Server.Host = "localhost"
	original.Server.Port = 8080
	original.Server.TLS.Enabled = true
	original.Server.TLS.Certificate = "/path/to/cert.pem"
	original.Server.TLS.Key = "/path/to/key.pem"
	original.Database.Driver = "postgres"
	original.Database.Host = "db.example.com"
	original.Database.Username = "admin"
	original.Database.Password = "password123"
	original.Database.Options = map[string]string{
		"sslmode": "require",
		"timeout": "30s",
	}
	original.Logging.Level = "info"
	original.Logging.Output = []string{"console", "file"}

	data, err := serializer.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	var result Config
	err = serializer.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	if !reflect.DeepEqual(original, result) {
		t.Errorf("Unmarshal() = %v, want %v", result, original)
	}
}
