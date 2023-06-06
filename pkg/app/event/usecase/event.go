package usecase

import (
	repo "nostr-relay/pkg/app/event/repo/postgre"
	"nostr-relay/pkg/models"
)

type Handler struct {
}

func NewEventHandler() *Handler {
	return &Handler{}
}

func (t *Handler) SaveEvent(data *models.RelayEvent) error {
	return repo.SaveEvent(data)
}

func (t *Handler) GetEvent(limit int) []models.RelayEvent {
	return repo.GetEvent(limit)
}

func (t *Handler) GetLastEvent() *models.RelayEvent {
	return repo.GetLastEvent()
}

func (t *Handler) GetEventFrom(id int) []models.RelayEvent {
	return repo.GetEventFrom(id)
}
