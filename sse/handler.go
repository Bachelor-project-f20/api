package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/99designs/gqlgen/graphql"
	gql "github.com/99designs/gqlgen/graphql"
)

type (
	SSE        struct{}
	connection struct {
		SSE
		ctx  context.Context
		exec gql.GraphExecutor
	}
	operationMessage struct {
		Payload json.RawMessage `json:"payload,omitempty"`
		ID      string          `json:"id,omitempty"`
		Type    string          `json:"type"`
	}
)

//type SSE struct {
// Upgrader              websocket.Upgrader
// InitFunc              WebsocketInitFunc
// KeepAlivePingInterval time.Duration
//}

func (t SSE) Supports(r *http.Request) bool {
	return true //r.Header.Get("Upgrade") != ""
}

func (t SSE) Do(w http.ResponseWriter, r *http.Request, exec graphql.GraphExecutor) {
	fmt.Println("Do!!!!!!!!!!!!!!!!!!")
	fmt.Println("Client: %v", r.RemoteAddr)
	fmt.Fprintf(w, "hello\n")
	// ws, err := t.Upgrader.Upgrade(w, r, http.Header{
	// 	"Sec-Websocket-Protocol": []string{"graphql-ws"},
	// })
	// if err != nil {
	// 	log.Printf("unable to upgrade %T to websocket %s: ", w, err.Error())
	// 	SendErrorf(w, http.StatusBadRequest, "unable to upgrade")
	// 	return
	// }

	// conn := connection{
	// 	//active:    map[string]context.CancelFunc{},
	// 	//conn:      ws,
	// 	ctx:  r.Context(),
	// 	exec: exec,
	// 	SSE:  t,
	// }

	fmt.Println("Request")
	fmt.Println(r)
	fmt.Println("Body")
	fmt.Println(r.Body)
	fmt.Println("---------------------------------")
	var p operationMessage

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Do something with the Person struct...
	fmt.Fprintf(w, "Person: %+v", p)

	// if !conn.init() {
	// 	return
	// }

	// conn.run()
}

// func (c *connection) readOp(r *http.Request) *operationMessage {
// 	// _, r, err := c.conn.NextReader()
// 	// if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseNoStatusReceived) {
// 	// 	return nil
// 	// } else if err != nil {
// 	// 	c.sendConnectionError("invalid json: %T %s", err, err.Error())
// 	// 	return nil
// 	// }
// 	// message := operationMessage{}
// 	// if err := jsonDecode(r.Body, &message); err != nil {
// 	// 	//c.sendConnectionError("invalid json")
// 	// 	return nil
// 	// }

// 	// return &message
// }
