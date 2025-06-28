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
	router.GET("/rocket/:id", controller.GetRocketByChannel)
	router.GET("/rockets", controller.GetRockets)

	router.Run("localhost:8088")
}
