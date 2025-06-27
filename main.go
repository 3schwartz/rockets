package main

import (
	"main/internal"

	"github.com/gin-gonic/gin"
)

func main() {
	storage := internal.NewMemoryStorage()
	controller := internal.Controller{Storage: &storage}

	router := gin.Default()

	router.POST("/messages", controller.PostMessages)
	router.GET("/messages/:id", controller.GetMessageByChannel)

	router.Run("localhost:8088")
}
