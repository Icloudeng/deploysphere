package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/domain", DomaineHandler)

	fmt.Println("Server running on PORT:", 8088)

	log.Fatal(r.Run(":8088"))
}
