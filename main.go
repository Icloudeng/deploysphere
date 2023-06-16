package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/provision", Handlers.provision)

	log.Println("Server running on PORT: ", 8088)
	log.Fatal(r.Run(":8088"))
}
