package controller

import (
	"encoding/json"
	"fmt"
	"main/internal/model"
	"main/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	Storage repository.IRepository
}

func (controller *Controller) PostMessages(c *gin.Context) {
	fmt.Println("postMessages")

	envelope, err := envelopeFromJson(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message, err := envelope.IntoRocketEvent()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	controller.Storage.Save(*message)

	jsonData, err := json.MarshalIndent(message, "", "  ")
	if err != nil {
		fmt.Println("JSON Marshal error:", err)
		return
	}
	fmt.Println(string(jsonData))
}

func (controller *Controller) GetMessageByChannel(c *gin.Context) {
	fmt.Println("getMessageByChannel")

	id := c.Param("id")

	events, err := controller.Storage.Events(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rocketState, err := model.RehydrateRocketState(events)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, rocketState)
}

func envelopeFromJson(c *gin.Context) (*model.Envelope, error) {
	var envelope *model.Envelope

	if err := c.ShouldBindJSON(&envelope); err != nil {
		return nil, err
	}

	return envelope, nil
}
