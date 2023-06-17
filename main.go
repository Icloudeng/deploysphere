package main

import (
	"flag"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	port := flag.Int("port", 8088, "Server port")
	init := flag.Bool("init", false, "Only initialize packackges")
	flag.Parse()

	if *init {
		return
	}

	r.POST("/provision", Handlers.provision)
	r.DELETE("/provision/:ref", Handlers.deleteProvision)

	log.Println("Server running on PORT: ", port)
	log.Fatal(r.Run(":" + strconv.Itoa(*port)))
}
