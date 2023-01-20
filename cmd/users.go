package main

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/joho/godotenv"
)

type User struct {
	FirstName string `fake:"{firstname}"`
	LastName  string `fake:"{lastname}"`
	Email     string `fake:"{email}"`
}

func main() {
	_ = godotenv.Load()

	var u User

	gofakeit.Struct(&u)

	println(u.FirstName)
}
