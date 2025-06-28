package internal

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	Storage IRepository
}

func (controller *Controller) PostMessages(c *gin.Context) {
	fmt.Println("postMessages")

	envelope, err := envelopeFromJson(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message, err := envelope.intoRocketEvent()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	controller.Storage.save(*message)

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

	events, err := controller.Storage.events(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rocketState, err := rehydrateRocketState(events)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, rocketState)
}

func envelopeFromJson(c *gin.Context) (*Envelope, error) {
	var envelope *Envelope

	if err := c.ShouldBindJSON(&envelope); err != nil {
		return nil, err
	}

	return envelope, nil
}
