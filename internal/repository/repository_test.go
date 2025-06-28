package repository

import (
	"main/internal/model"
	"sort"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestWhenSaveThenAppleToFecth(t *testing.T) {
	// Arrange
	repo := NewMemoryStorage()
	channel := "a"
	rocketEvent := model.RocketEvent{
		Channel:       channel,
		MessageNumber: 3,
		MessageTime:   time.Now(),
		Message:       model.RocketExploded{Reason: "some reason"},
	}

	// Act
	repo.Save(rocketEvent)
	events, err := repo.EventsByChannel(channel)

	// Assert
	if err != nil {
		t.Error(err)
	}
	if len(events) != 1 {
		t.Errorf("expected length of 1 got %d", len(events))
	}
}

func TestGivenNoMatchinChannelThenReturnError(t *testing.T) {
	// Arrange
	repo := NewMemoryStorage()

	// Act
	_, err := repo.EventsByChannel("")

	// Assert
	if err == nil {
		t.Error("expected error")
	}
}

func TestWhenFilterExplodedThenReturnOnlyExplodedAndLimit(t *testing.T) {
	// Arrange
	repo := NewMemoryStorage()
	now := time.Now()
	repo.Save(model.RocketEvent{
		Channel:       "a",
		MessageNumber: 3,
		MessageTime:   now,
		Message:       model.RocketExploded{Reason: "some reason"},
	})
	repo.Save(model.RocketEvent{
		Channel:       "b",
		MessageNumber: 3,
		MessageTime:   now,
		Message:       model.RocketSpeedDecreased{By: 1},
	})
	repo.Save(model.RocketEvent{
		Channel:       "c",
		MessageNumber: 3,
		MessageTime:   now,
		Message:       model.RocketSpeedDecreased{By: 1},
	})

	// Act
	events := repo.Events(true, -1)

	// Assert
	if len(events) != 2 {
		t.Errorf("expected length of 2 got %d", len(events))
	}
	expected := [][]model.RocketEvent{
		{{
			Channel:       "b",
			MessageNumber: 3,
			MessageTime:   now,
			Message:       model.RocketSpeedDecreased{By: 1},
		}},
		{{
			Channel:       "c",
			MessageNumber: 3,
			MessageTime:   now,
			Message:       model.RocketSpeedDecreased{By: 1},
		}},
	}
	sortByChannel(&expected)
	sortByChannel(&events)
	if diff := cmp.Diff(expected, events); diff != "" {
		t.Error(diff)
	}
}

func TestWhenFilterExplodedAndLimitThenReturnOnlyExplodedAndLimit(t *testing.T) {
	// Arrange
	repo := NewMemoryStorage()
	now := time.Now()
	repo.Save(model.RocketEvent{
		Channel:       "a",
		MessageNumber: 3,
		MessageTime:   now,
		Message:       model.RocketExploded{Reason: "some reason"},
	})
	repo.Save(model.RocketEvent{
		Channel:       "b",
		MessageNumber: 3,
		MessageTime:   now,
		Message:       model.RocketSpeedDecreased{By: 1},
	})
	repo.Save(model.RocketEvent{
		Channel:       "c",
		MessageNumber: 3,
		MessageTime:   now,
		Message:       model.RocketSpeedDecreased{By: 1},
	})

	// Act
	events := repo.Events(true, 1)

	// Assert
	if len(events) != 1 {
		t.Errorf("expected length of 1 got %d", len(events))
	}
	expected := [][]model.RocketEvent{
		{{
			Channel:       "b",
			MessageNumber: 3,
			MessageTime:   now,
			Message:       model.RocketSpeedDecreased{By: 1},
		}},
	}
	sortByChannel(&expected)
	sortByChannel(&events)
	if diff := cmp.Diff(expected, events); diff != "" {
		t.Error(diff)
	}
}

func TestWhenLimitThenReturnOnlLimit(t *testing.T) {
	// Arrange
	repo := NewMemoryStorage()
	now := time.Now()
	repo.Save(model.RocketEvent{
		Channel:       "a",
		MessageNumber: 3,
		MessageTime:   now,
		Message:       model.RocketExploded{Reason: "some reason"},
	})
	repo.Save(model.RocketEvent{
		Channel:       "b",
		MessageNumber: 3,
		MessageTime:   now,
		Message:       model.RocketSpeedDecreased{By: 1},
	})
	repo.Save(model.RocketEvent{
		Channel:       "c",
		MessageNumber: 3,
		MessageTime:   now,
		Message:       model.RocketSpeedDecreased{By: 1},
	})

	// Act
	events := repo.Events(false, 2)

	// Assert
	if len(events) != 2 {
		t.Errorf("expected length of 2 got %d", len(events))
	}
}

func sortByChannel(events *[][]model.RocketEvent) {
	sort.Slice(*events, func(i, j int) bool {
		return (*events)[i][0].Channel < (*events)[j][0].Channel
	})
}
