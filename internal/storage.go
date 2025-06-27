package internal

import "fmt"

type Storage interface {
	save(message RocketEvent)
	events(channel string) ([]RocketEvent, error)
}

type MemoryStorage struct {
	messages map[string][]RocketEvent
}

func NewMemoryStorage() MemoryStorage {
	return MemoryStorage{
		messages: map[string][]RocketEvent{},
	}
}

func (s *MemoryStorage) save(message RocketEvent) {
	channel := message.Channel
	s.messages[channel] = append(s.messages[channel], message)
}

func (s *MemoryStorage) events(channel string) ([]RocketEvent, error) {
	events, ok := s.messages[channel]
	if !ok {
		return nil, fmt.Errorf("channel %q not found", channel)
	}
	return events, nil
}
