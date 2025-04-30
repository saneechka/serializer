package json

import (
	"reflect"
	"testing"
)

func TestJSONSerializer(t *testing.T) {
	type TestStruct struct {
		String  string   `json:"string"`
		Integer int      `json:"integer"`
		Float   float64  `json:"float"`
		Boolean bool     `json:"boolean"`
		Array   []string `json:"array"`
		Nested  struct {
			Field string `json:"field"`
		} `json:"nested"`
	}

	serializer := New()

	if format := serializer.Format(); format != "JSON" {
		t.Errorf("Format() = %v, want %v", format, "JSON")
	}

	original := TestStruct{
		String:  "тест",
		Integer: 42,
		Float:   3.14,
		Boolean: true,
		Array:   []string{"один", "два", "три"},
		Nested: struct {
			Field string `json:"field"`
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

	err = serializer.Unmarshal([]byte("{invalid json}"), &result)
	if err == nil {
		t.Error("Unmarshal() с неверным JSON должен возвращать ошибку")
	}
}

func TestJSONNilValues(t *testing.T) {
	serializer := New()

	data, err := serializer.Marshal(nil)
	if err != nil {
		t.Errorf("Marshal(nil) error = %v", err)
	}

	if string(data) != "null" {
		t.Errorf("Marshal(nil) = %v, want %v", string(data), "null")
	}

	var result interface{}
	err = serializer.Unmarshal([]byte("null"), &result)
	if err != nil {
		t.Errorf("Unmarshal('null') error = %v", err)
	}
}
