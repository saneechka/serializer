package json

import (
	"encoding/json"
	"reflect"
	"testing"
)

type TestStruct struct {
	Name    string     `json:"name"`
	Age     int        `json:"age"`
	Hobbies []string   `json:"hobbies"`
	Info    InfoStruct `json:"info"`
}

type InfoStruct struct {
	Address string `json:"address"`
	Email   string `json:"email"`
}

func TestMarshal(t *testing.T) {
	testObj := TestStruct{
		Name: "Алексей",
		Age:  30,
		Hobbies: []string{
			"чтение",
			"программирование",
		},
		Info: InfoStruct{
			Address: "Москва",
			Email:   "alex@example.com",
		},
	}

	serializer := New()

	got, err := serializer.Marshal(testObj)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	want, err := json.Marshal(testObj)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Marshal() = %s, want %s", string(got), string(want))
	}
}

func TestUnmarshal(t *testing.T) {
	expected := TestStruct{
		Name: "Алексей",
		Age:  30,
		Hobbies: []string{
			"чтение",
			"программирование",
		},
		Info: InfoStruct{
			Address: "Москва",
			Email:   "alex@example.com",
		},
	}

	jsonData := `{"name":"Алексей","age":30,"hobbies":["чтение","программирование"],"info":{"address":"Москва","email":"alex@example.com"}}`

	serializer := New()

	var got TestStruct
	if err := serializer.Unmarshal([]byte(jsonData), &got); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Unmarshal() got = %v, want %v", got, expected)
	}
}

func TestFormat(t *testing.T) {
	serializer := New()

	if got := serializer.Format(); got != "JSON" {
		t.Errorf("Format() = %v, want %v", got, "JSON")
	}
}

func TestUnmarshalError(t *testing.T) {
	invalidJSON := []byte(`{"name": "Алексей", "age": "тридцать"}`)

	serializer := New()

	var got TestStruct
	if err := serializer.Unmarshal(invalidJSON, &got); err == nil {
		t.Error("Unmarshal() expected error with invalid JSON, got nil")
	}
}

func TestComplexStructures(t *testing.T) {
	jsonData := `{
		"name": "Проект",
		"version": "1.0.0",
		"developers": [
			{"name": "Алексей", "role": "Lead"},
			{"name": "Мария", "role": "Backend"}
		],
		"config": {
			"debug": true,
			"database": {
				"host": "localhost",
				"port": 5432
			}
		}
	}`

	type ComplexStruct struct {
		Name       string `json:"name"`
		Version    string `json:"version"`
		Developers []struct {
			Name string `json:"name"`
			Role string `json:"role"`
		} `json:"developers"`
		Config struct {
			Debug    bool `json:"debug"`
			Database struct {
				Host string `json:"host"`
				Port int    `json:"port"`
			} `json:"database"`
		} `json:"config"`
	}

	serializer := New()

	var got ComplexStruct
	if err := serializer.Unmarshal([]byte(jsonData), &got); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	if got.Name != "Проект" {
		t.Errorf("Name = %v, want %v", got.Name, "Проект")
	}
	if len(got.Developers) != 2 {
		t.Errorf("len(Developers) = %v, want %v", len(got.Developers), 2)
	}
	if got.Config.Database.Port != 5432 {
		t.Errorf("Config.Database.Port = %v, want %v", got.Config.Database.Port, 5432)
	}
}

func TestNullValues(t *testing.T) {
	jsonData := `{"name":"Алексей","age":30,"hobbies":null,"info":null}`

	serializer := New()

	var got TestStruct
	if err := serializer.Unmarshal([]byte(jsonData), &got); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	if got.Name != "Алексей" {
		t.Errorf("Name = %v, want %v", got.Name, "Алексей")
	}
	if got.Age != 30 {
		t.Errorf("Age = %v, want %v", got.Age, 30)
	}

	if got.Hobbies != nil && len(got.Hobbies) > 0 {
		t.Errorf("Hobbies = %v, want nil or empty slice", got.Hobbies)
	}

	if got.Info.Address != "" || got.Info.Email != "" {
		t.Errorf("Info = %v, want zero values", got.Info)
	}
}
