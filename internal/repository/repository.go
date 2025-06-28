package repository

import (
	"fmt"
	"main/internal/model"
)

type IRepository interface {
	Save(message model.RocketEvent)
	Events(channel string) ([]model.RocketEvent, error)
}

type MemoryRepository struct {
	messages map[string][]model.RocketEvent
}

func NewMemoryStorage() MemoryRepository {
	return MemoryRepository{
		messages: map[string][]model.RocketEvent{},
	}
}

func (s *MemoryRepository) Save(message model.RocketEvent) {
	channel := message.Channel
	s.messages[channel] = append(s.messages[channel], message)
}

func (s *MemoryRepository) Events(channel string) ([]model.RocketEvent, error) {
	events, ok := s.messages[channel]
	if !ok {
		return nil, fmt.Errorf("channel %q not found", channel)
	}
	return events, nil
}
