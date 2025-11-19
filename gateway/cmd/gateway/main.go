package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("main: gateway start")
	if err := run(); err != nil {
		log.Fatal("gateway main: " + err.Error())
	}
	fmt.Println("main: gateway start end")
}

func run() error {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Run(":8080")
	return nil
}
