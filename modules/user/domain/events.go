package domain

import "time"

type EventUserRegistered struct {
	Email      string
	OccurredAt time.Time
}

func NewEventUserRegistered(email string) *EventUserRegistered {
	return &EventUserRegistered{
		Email:      email,
		OccurredAt: time.Now(),
	}
}
