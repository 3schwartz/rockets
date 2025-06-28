package model

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"
)

type Envelope struct {
	Message  json.RawMessage `json:"message"` // Keep raw until you know type
	Metadata Metadata        `json:"metadata"`
}

func (e *Envelope) IntoRocketEvent() (*RocketEvent, error) {
	var message RocketMessage

	switch e.Metadata.MessageType {
	case "RocketSpeedDecreased":
		var msg RocketSpeedDecreased
		if err := json.Unmarshal(e.Message, &msg); err != nil {
			return nil, err
		}
		fmt.Printf("Rocket speed decreased by %d\n", msg.By)
		message = msg

	case "RocketSpeedIncreased":
		var msg RocketSpeedIncreased
		if err := json.Unmarshal(e.Message, &msg); err != nil {
			return nil, err
		}
		fmt.Printf("Rocket speed increased by %d\n", msg.By)
		message = msg

	case "RocketExploded":
		var msg RocketExploded
		if err := json.Unmarshal(e.Message, &msg); err != nil {
			return nil, err
		}
		fmt.Printf("Rocket exploded due to: %s\n", msg.Reason)
		message = msg

	case "RocketLaunched":
		var msg RocketLaunched
		if err := json.Unmarshal(e.Message, &msg); err != nil {
			return nil, err
		}
		fmt.Printf("Rocket launched: type=%s, speed=%d, mission=%s\n", msg.Type, msg.LaunchSpeed, msg.Mission)
		message = msg

	case "RocketMissionChanged":
		var msg RocketMissionChanged
		if err := json.Unmarshal(e.Message, &msg); err != nil {
			return nil, err
		}
		fmt.Printf("Rocket mission changed to: %s\n", msg.NewMission)
		message = msg

	default:
		return nil, fmt.Errorf("unsupported event: %s", e.Metadata.MessageType)
	}

	return &RocketEvent{
		Channel:       e.Metadata.Channel,
		MessageNumber: e.Metadata.MessageNumber,
		MessageTime:   e.Metadata.MessageTime,
		Message:       message,
	}, nil
}

type Metadata struct {
	Channel       string    `json:"channel"`
	MessageNumber int       `json:"messageNumber"`
	MessageTime   time.Time `json:"messageTime"`
	MessageType   string    `json:"messageType"`
}

type RocketMessage interface {
	eventName() string
}

type RocketLaunched struct {
	Type        string `json:"type"`
	LaunchSpeed int    `json:"launchSpeed"`
	Mission     string `json:"mission"`
}

func (r RocketLaunched) eventName() string {
	return "RocketLaunched"
}

type RocketSpeedDecreased struct {
	By int `json:"by"`
}

func (r RocketSpeedDecreased) eventName() string {
	return "RocketSpeedDecreased"
}

type RocketSpeedIncreased struct {
	By int `json:"by"`
}

func (r RocketSpeedIncreased) eventName() string {
	return "RocketSpeedIncreased"
}

type RocketMissionChanged struct {
	NewMission string `json:"newMission"`
}

func (r RocketMissionChanged) eventName() string {
	return "RocketMissionChanged"
}

type RocketExploded struct {
	Reason string `json:"reason"`
}

func (r RocketExploded) eventName() string {
	return "RocketExploded"
}

type RocketEvent struct {
	Channel       string
	MessageNumber int
	MessageTime   time.Time
	Message       RocketMessage
}

type RocketState struct {
	Type     string `json:"type"`
	Speed    int    `json:"speed"`
	Mission  string `json:"mission"`
	Exploded bool   `json:"exploded"`
}

func RehydrateRocketState(events []RocketEvent) (*RocketState, error) {
	rocket := &RocketState{}

	if len(events) == 0 {
		return nil, fmt.Errorf("some events needs to be in place")
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].MessageNumber < events[j].MessageNumber
	})

	firstEvent := events[0]
	if _, ok := firstEvent.Message.(RocketLaunched); !ok {
		return nil, fmt.Errorf("need to have initial event to construct state")
	}

	for _, event := range events {
		switch e := event.Message.(type) {
		case RocketLaunched:
			rocket.Type = e.Type
			rocket.Speed = e.LaunchSpeed
			rocket.Mission = e.Mission
		case RocketSpeedIncreased:
			rocket.Speed += e.By
		case RocketSpeedDecreased:
			rocket.Speed -= e.By
		case RocketMissionChanged:
			rocket.Mission = e.NewMission
		case RocketExploded:
			rocket.Exploded = true
		default:
			return nil, fmt.Errorf("unknown state: %s", e.eventName())
		}
	}

	return rocket, nil
}
