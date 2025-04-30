package toml

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/BurntSushi/toml"
)


type TestStruct struct {
	Name    string     `toml:"name"`
	Age     int        `toml:"age"`
	Hobbies []string   `toml:"hobbies"`
	Info    InfoStruct `toml:"info"`
}


type InfoStruct struct {
	Address string `toml:"address"`
	Email   string `toml:"email"`
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


	var buf bytes.Buffer
	encoder := toml.NewEncoder(&buf)
	if err := encoder.Encode(testObj); err != nil {
		t.Fatalf("toml.NewEncoder().Encode() error = %v", err)
	}
	want := buf.Bytes()


	var gotMap, wantMap map[string]interface{}
	if err := toml.Unmarshal(got, &gotMap); err != nil {
		t.Fatalf("toml.Unmarshal(got) error = %v", err)
	}
	if err := toml.Unmarshal(want, &wantMap); err != nil {
		t.Fatalf("toml.Unmarshal(want) error = %v", err)
	}

	if !reflect.DeepEqual(gotMap, wantMap) {
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


	tomlData := `
name = "Алексей"
age = 30
hobbies = ["чтение", "программирование"]

[info]
address = "Москва"
email = "alex@example.com"
`


	serializer := New()


	var got TestStruct
	if err := serializer.Unmarshal([]byte(tomlData), &got); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}


	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Unmarshal() got = %v, want %v", got, expected)
	}
}

func TestFormat(t *testing.T) {

	serializer := New()


	if got := serializer.Format(); got != "TOML" {
		t.Errorf("Format() = %v, want %v", got, "TOML")
	}
}


func TestUnmarshalError(t *testing.T) {

	invalidTOML := []byte(`
name = "Алексей"
age = "тридцать" # должно быть число, а не строка
`)


	serializer := New()


	var got TestStruct
	if err := serializer.Unmarshal(invalidTOML, &got); err == nil {
		t.Error("Unmarshal() expected error with invalid TOML, got nil")
	}
}


func TestComplexStructures(t *testing.T) {

	tomlData := `
name = "Проект"
version = "1.0.0"

[[developers]]
name = "Алексей"
role = "Lead"

[[developers]]
name = "Мария"
role = "Backend"

[config]
debug = true

[config.database]
host = "localhost"
port = 5432
`


	type ComplexStruct struct {
		Name       string `toml:"name"`
		Version    string `toml:"version"`
		Developers []struct {
			Name string `toml:"name"`
			Role string `toml:"role"`
		} `toml:"developers"`
		Config struct {
			Debug    bool `toml:"debug"`
			Database struct {
				Host string `toml:"host"`
				Port int    `toml:"port"`
			} `toml:"database"`
		} `toml:"config"`
	}


	serializer := New()


	var got ComplexStruct
	if err := serializer.Unmarshal([]byte(tomlData), &got); err != nil {
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
