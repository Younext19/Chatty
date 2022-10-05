package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("%s", err.Error())
	}

}
