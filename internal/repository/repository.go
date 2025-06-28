package repository

import (
	"fmt"
	"main/internal/model"
)

// IRepository defines methods to persist and query Rocket events.
type IRepository interface {
	// Save stores a RocketEvent.
	Save(message model.RocketEvent)

	// EventsByChannel returns all events for a specific channel.
	EventsByChannel(channel string) ([]model.RocketEvent, error)

	// Events returns all events with optional filters.
	//
	// If filterExploded is true, exploded events are excluded.
	// If limit >= 0, only up to `limit` events are returned per channel.
	Events(filterExploded bool, limit int) [][]model.RocketEvent
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

func (s *MemoryRepository) EventsByChannel(channel string) ([]model.RocketEvent, error) {
	events, ok := s.messages[channel]
	if !ok {
		return nil, fmt.Errorf("channel %q not found", channel)
	}
	return events, nil
}

func (s *MemoryRepository) Events(filterExploded bool, limit int) [][]model.RocketEvent {
	channels := [][]model.RocketEvent{}

	event_count := 0
	for _, events := range s.messages {
		if limit >= 0 && event_count >= limit {
			break
		}
		if filterExploded {
			isExploded := false
			for _, event := range events {
				if _, ok := event.Message.(model.RocketExploded); ok {
					isExploded = true
					break
				}
			}
			if isExploded {
				continue
			}
		}
		channels = append(channels, events)
		event_count += 1
	}

	return channels
}
