package sse

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	ob "github.com/Bachelor-project-f20/go-outbox"
	models "github.com/Bachelor-project-f20/shared/models"
	"github.com/golang/protobuf/proto"
)

type Client struct {
	ID string
}

type SSEHandler struct {
	outbox         ob.Outbox
	eventChan      <-chan models.Event
	clients        map[string]chan string
	closingClients chan string
}

func NewSSEHandler(outbox ob.Outbox, eventChan <-chan models.Event) SSEHandler {
	sseh := SSEHandler{
		outbox,
		eventChan,
		make(map[string]chan string),
		make(chan string, 10),
	}

	go sseh.listen()

	return sseh
}

func payloadToJson(event models.Event) string {
	switch event.EventName {
	case models.UserEvents_USER_CREATED.String():
		var m models.UserCreated
		_ = proto.Unmarshal(event.Payload, &m)
		b, _ := json.Marshal(m)
		return string(b)
	case models.UserEvents_USER_UPDATED.String():
		var m models.UserUpdated
		_ = proto.Unmarshal(event.Payload, &m)
		b, _ := json.Marshal(m)
		return string(b)
	case models.UserEvents_USER_DELETED.String():
		var m models.UserDeleted
		_ = proto.Unmarshal(event.Payload, &m)
		b, _ := json.Marshal(m)
		return string(b)
	}
	return "error"
}

func (s *SSEHandler) listen() {
	for {
		select {
		case name := <-s.closingClients:
			log.Println("Close connection to client: ", name)
			delete(s.clients, name)
		case event := <-s.eventChan:
			fmt.Println("Recived event: ", event.ID)
			msgChan := s.clients[event.ApiTag]
			fmt.Println("msgChan ", msgChan != nil, " name: ", event.ApiTag)
			if msgChan != nil {
				msgChan <- payloadToJson(event)
			}
		}
	}
}

func (s *SSEHandler) Handler(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}
	//log.Printf("Client: %v \n", r.RemoteAddr)
	msgChan := s.clients[r.RemoteAddr]
	if msgChan == nil {
		log.Println("New client: ", r.RemoteAddr)
		msgChan = make(chan string, 10)
		s.clients[r.RemoteAddr] = msgChan
		msgChan <- "{\"apiTag\": " + "\"" + r.RemoteAddr + "\"}"
	}
	//var client Client
	//s.outbox.GetDBConnection().FirstOrCreate(&client, Client{ID: r.RemoteAddr})

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	defer func() {
		s.closingClients <- r.RemoteAddr
	}()

	//fmt.Fprintf(w, "data: %v\n\n", "ping")
	//fmt.Fprint(w, "data: ", rand.Intn(100), ",", rand.Intn(100), "\n\n")
	//fmt.Printf("HERE")
	//fmt.Fprintf(w, ": no\n\n")

	notify := r.Context().Done()
	go func() {
		<-notify
		s.closingClients <- r.RemoteAddr
	}()

	for {
		log.Println("Start sending")
		fmt.Fprintf(w, "data: %v\n\n", <-msgChan)
		fmt.Fprintf(w, ": no\n\n")
		flusher.Flush()
		log.Println("Send done")
	}

	// timeout := time.After(10 * time.Second)
	// select {
	// case ev := <-s.eventChan:
	// 	var buf bytes.Buffer
	// 	enc := json.NewEncoder(&buf)
	// 	enc.Encode(ev)
	// 	fmt.Fprintf(w, "data: %v\n\n", buf.String())
	// 	fmt.Printf("data: %v\n", buf.String())
	// case <-timeout:
	// 	fmt.Println("Time out")
	// 	fmt.Fprintf(w, "event: %v\n", "ping")
	// 	fmt.Fprint(w, "data: ", rand.Intn(100), ",", rand.Intn(100), "\n\n")
	// 	// fmt.Fprintf(w, ": no\n\n")
	// }

	//Send event
	// if f, ok := w.(http.Flusher); ok {
	// 	f.Flush()
	// }
}
