package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/provision", Handlers.provision)
	r.DELETE("/provision/:ref", Handlers.deleteProvision)

	log.Println("Server running on PORT: ", 8088)
	log.Fatal(r.Run(":8088"))
}
