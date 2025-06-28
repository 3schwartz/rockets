package model

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestGivenEventsThenConstructState(t *testing.T) {
	// Arrange
	channel := "a"
	events := []RocketEvent{
		{
			Channel:       channel,
			MessageNumber: 1,
			MessageTime:   time.Now(),
			Message:       RocketLaunched{Type: "hello", LaunchSpeed: 10, Mission: "world"},
		},
		{
			Channel:       channel,
			MessageNumber: 3,
			MessageTime:   time.Now(),
			Message:       RocketSpeedDecreased{By: 2},
		},
		{
			Channel:       channel,
			MessageNumber: 2,
			MessageTime:   time.Now(),
			Message:       RocketSpeedIncreased{By: 1},
		},
	}

	// Act
	rocketState, err := RehydrateRocketState(events)

	// Assert
	expected := RocketState{
		Type:     "hello",
		Speed:    9,
		Mission:  "world",
		Exploded: false,
	}
	if err != nil {
		t.Error(err)
	}
	if diff := cmp.Diff(&expected, rocketState); diff != "" {
		t.Error(diff)
	}
}
