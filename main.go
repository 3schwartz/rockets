package main

import (
	"main/internal/controller"
	"main/internal/repository"

	"github.com/gin-gonic/gin"
)

func main() {
	storage := repository.NewMemoryStorage()
	controller := controller.Controller{Storage: &storage}

	router := gin.Default()

	router.POST("/messages", controller.PostMessages)
	router.GET("/messages/:id", controller.GetMessageByChannel)

	router.Run("localhost:8088")
}
