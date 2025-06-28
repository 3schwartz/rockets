package controller

import (
	"encoding/json"
	"fmt"
	"main/internal/model"
	"main/internal/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	Storage repository.IRepository
}

// Post rocket messages
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

// GetRocketByChannel returns the reconstructed rocket state for a given channel ID.
//
// It looks up all events associated with the provided channel ID,
// rehydrates the rocket state by replaying them, and returns the result.
//
// Example:
//
//	curl http://localhost:8088/rocket/<channel-id>
func (controller *Controller) GetRocketByChannel(c *gin.Context) {
	fmt.Println("getRocketByChannel")

	id := c.Param("id")

	events, err := controller.Storage.EventsByChannel(id)
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

// GetRockets returns a list of rocket states across all channels.
//
// Optional query parameters:
//   - filterExploded (bool): If true, only non-exploded rockets are included.
//   - limit (int): Maximum number of events per channel. Use -1 for no limit.
//
// Example:
//
//	curl http://localhost:8088/rockets?filterExploded=true&limit=2
func (controller *Controller) GetRockets(c *gin.Context) {
	filterExplodedStr := c.DefaultQuery("filterExploded", "false")
	filterExploded, err := strconv.ParseBool(filterExplodedStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid value for filterExploded"})
		return
	}

	limitStr := c.DefaultQuery("limit", "-1")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid value for limit"})
		return
	}

	events := controller.Storage.Events(filterExploded, limit)

	rocketStates := []model.RocketState{}
	for _, rocketEvents := range events {
		rocketState, err := model.RehydrateRocketState(rocketEvents)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		rocketStates = append(rocketStates, *rocketState)
	}

	c.IndentedJSON(http.StatusOK, rocketStates)
}

func envelopeFromJson(c *gin.Context) (*model.Envelope, error) {
	var envelope *model.Envelope

	if err := c.ShouldBindJSON(&envelope); err != nil {
		return nil, err
	}

	return envelope, nil
}
