package internal

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	Storage Storage
}

func (controller *Controller) PostMessages(c *gin.Context) {
	fmt.Println("postMessages")

	envelope, err := EnvelopeFromJson(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message, err := envelope.IntoRocketEvent()
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
