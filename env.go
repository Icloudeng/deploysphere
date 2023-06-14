package main

import (
	"log"

	"github.com/joho/godotenv"
)

var (
	Token string
)

func init() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
