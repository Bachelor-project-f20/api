package graphql

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

import (
	"context"
	"log"
	"time"

	"github.com/Bachelor-project-f20/eventToGo"
	models "github.com/Bachelor-project-f20/shared/models"
	"github.com/gofrs/uuid"
	"github.com/golang/protobuf/proto"
)

const pubAPI string = "api"

type Resolver struct {
	Emitter eventToGo.EventEmitter
}

func (r *mutationResolver) emitEvent(payLoad []byte, eventType models.UserEvents) (bool, error) {
	id, _ := uuid.NewV4()

	event := &models.Event{
		ID:        id.String(),
		Publisher: pubAPI,
		EventName: eventType.String(),
		Timestamp: time.Now().UnixNano(),
		Payload:   payLoad,
	}
	r.Emitter.Emit(*event)
	return true, nil
}

func (r *mutationResolver) CreateUser(ctx context.Context, user models.CreateUser) (bool, error) {
	marshalEvent, err := proto.Marshal(&user)
	if err != nil {
		log.Fatalf("Error with proto: %v \n", err)
		return false, err
	}

	return r.emitEvent(marshalEvent, models.UserEvents_CREATE_USER)
}

func (r *mutationResolver) UpdateUser(ctx context.Context, user models.UpdateUser) (bool, error) {
	marshalEvent, err := proto.Marshal(&user)
	if err != nil {
		log.Fatalf("Error with proto: %v \n", err)
		return false, err
	}

	return r.emitEvent(marshalEvent, models.UserEvents_UPDATE_USER)
}

func (r *mutationResolver) DeleteUser(ctx context.Context, user models.DeleteUser) (bool, error) {
	marshalEvent, err := proto.Marshal(&user)
	if err != nil {
		log.Fatalf("Error with proto: %v \n", err)
		return false, err
	}

	return r.emitEvent(marshalEvent, models.UserEvents_DELETE_USER)
}

//DUMMY func
func (r *queryResolver) GetUser(ctx context.Context, user models.User) (bool, error) {
	return false, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
