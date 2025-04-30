package main

import (
	"fmt"
	"log"

	"github.com/saneechka/serializer"
	"github.com/saneechka/serializer/json"
	"github.com/saneechka/serializer/toml"
)


type Person struct {
	Name    string  `json:"name" toml:"name"`
	Age     int     `json:"age" toml:"age"`
	Address Address `json:"address" toml:"address"`
}


type Address struct {
	City    string `json:"city" toml:"city"`
	Country string `json:"country" toml:"country"`
}

func main() {

	person := Person{
		Name: "Алексей",
		Age:  30,
		Address: Address{
			City:    "Москва",
			Country: "Россия",
		},
	}


	jsonSerializer := json.New()
	tomlSerializer := toml.New()


	fmt.Printf("Сериализация в %s:\n", jsonSerializer.Format())
	demonstrateSerializer(jsonSerializer, person)


	fmt.Printf("\nСериализация в %s:\n", tomlSerializer.Format())
	demonstrateSerializer(tomlSerializer, person)
}


func demonstrateSerializer(s serializer.Serializer, person Person) {
	// Сериализация
	data, err := s.Marshal(person)
	if err != nil {
		log.Fatalf("Ошибка сериализации: %v", err)
	}
	fmt.Printf("Сериализованные данные:\n%s\n", string(data))

	// Десериализация
	var newPerson Person
	if err := s.Unmarshal(data, &newPerson); err != nil {
		log.Fatalf("Ошибка десериализации: %v", err)
	}

	fmt.Printf("Десериализованный объект: %+v\n", newPerson)
}
