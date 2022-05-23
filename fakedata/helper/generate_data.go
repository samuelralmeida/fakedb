package helper

import (
	"fmt"

	"github.com/brianvoe/gofakeit/v6"
)

func GenerateAddress() string {
	types := []string{"Avenida", "Rua", "Travessa"}

	street := fmt.Sprintf("%s %s %s, %d",
		types[gofakeit.Number(0, 2)],
		gofakeit.StreetPrefix(),
		gofakeit.StreetName(),
		gofakeit.Number(0, 3000),
	)
	return street
}

func GenerateName() string {
	return fmt.Sprintf("%s %s %s",
		gofakeit.FirstName(),
		gofakeit.LastName(),
		gofakeit.LastName(),
	)
}
