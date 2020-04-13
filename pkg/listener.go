package api

import (
	"github.com/Bachelor-project-f20/eventToGo"
	models "github.com/Bachelor-project-f20/shared/models"
)

func setupEventListener(eventListener eventToGo.EventListener) (<-chan models.Event, error) {
	incomingEvents := []string{
		//
		//User events
		//
		models.UserEvents_USER_CREATED.String(),
		models.UserEvents_USER_DELETED.String(),
		models.UserEvents_USER_UPDATED.String()}

	eventChan, _, err := eventListener.Listen(incomingEvents...)
	return eventChan, err
}
