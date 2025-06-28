package internal

import "fmt"

type IRepository interface {
	save(message RocketEvent)
	events(channel string) ([]RocketEvent, error)
}

type MemoryRepository struct {
	messages map[string][]RocketEvent
}

func NewMemoryStorage() MemoryRepository {
	return MemoryRepository{
		messages: map[string][]RocketEvent{},
	}
}

func (s *MemoryRepository) save(message RocketEvent) {
	channel := message.Channel
	s.messages[channel] = append(s.messages[channel], message)
}

func (s *MemoryRepository) events(channel string) ([]RocketEvent, error) {
	events, ok := s.messages[channel]
	if !ok {
		return nil, fmt.Errorf("channel %q not found", channel)
	}
	return events, nil
}
