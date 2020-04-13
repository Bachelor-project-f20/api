package graphql

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

import (
	"context"
	"fmt"
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

func (r *mutationResolver) emitEvent(payLoad []byte, eventType models.UserEvents, apiTag string) (bool, error) {
	id, _ := uuid.NewV4()

	event := &models.Event{
		ID:        id.String(),
		Publisher: pubAPI,
		EventName: eventType.String(),
		Timestamp: time.Now().UnixNano(),
		Payload:   payLoad,
		ApiTag:    apiTag,
	}
	r.Emitter.Emit(*event)
	return true, nil
}

func (r *mutationResolver) CreateUser(ctx context.Context, user models.CreateUser, id string) (bool, error) {
	marshalEvent, err := proto.Marshal(&user)
	if err != nil {
		log.Fatalf("Error with proto: %v \n", err)
		return false, err
	}

	return r.emitEvent(marshalEvent, models.UserEvents_CREATE_USER, id)
}

func (r *mutationResolver) UpdateUser(ctx context.Context, user models.UpdateUser, id string) (bool, error) {
	marshalEvent, err := proto.Marshal(&user)
	if err != nil {
		log.Fatalf("Error with proto: %v \n", err)
		return false, err
	}

	return r.emitEvent(marshalEvent, models.UserEvents_UPDATE_USER, id)
}

func (r *mutationResolver) DeleteUser(ctx context.Context, user models.DeleteUser, id string) (bool, error) {
	marshalEvent, err := proto.Marshal(&user)
	if err != nil {
		log.Fatalf("Error with proto: %v \n", err)
		return false, err
	}

	return r.emitEvent(marshalEvent, models.UserEvents_DELETE_USER, id)
}

//DUMMY func
func (r *queryResolver) GetUser(ctx context.Context, user models.User) (bool, error) {
	return false, nil
}

func (r *subscriptionResolver) UserJoined(ctx context.Context, user string) (<-chan string, error) {
	fmt.Println("USERJOINED")
	c := make(chan string)
	go func() {
		for {
			time.Sleep(1000 * time.Millisecond)
			c <- "test"
		}

	}()
	return c, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Subscription returns SubscriptionResolver implementation.
func (r *Resolver) Subscription() SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
