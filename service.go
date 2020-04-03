package api

import (
	"log"
	"net/http"

	"github.com/Bachelor-project-f20/eventToGo"

	"github.com/99designs/gqlgen/handler"
	"github.com/Bachelor-project-f20/api/graphql"
	"github.com/Bachelor-project-f20/shared/config"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
)

type service struct {
	emitter eventToGo.EventEmitter
}

func Run() {
	configFile := "###" //TODO
	configRes, err := config.ConfigService(
		configFile,
		config.ConfigValues{
			UseEmitter: true,
		},
	)
	if err != nil {
		log.Fatalln("configuration failed, error: ", err)
		panic("configuration failed")
	}

	resolver := graphql.Resolver{
		Emitter: configRes.EventEmitter,
	}

	router := chi.NewRouter()

	// Add CORS middleware around every request
	// See https://github.com/rs/cors for full option listing
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"X-Requested-With", "Accept", "Authorization", "Accept-Language", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}).Handler)

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			// Check against your desired domains here
			return r.Host == "http://localhost:8081" || r.Host == "http://localhost:3000"
		},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	router.Handle("/", handler.Playground("GraphQL Playground", "/api"))
	router.Handle("/api",
		handler.GraphQL(graphql.NewExecutableSchema(graphql.Config{Resolvers: &resolver}), handler.WebsocketUpgrader(upgrader)),
	)

	log.Println("API: Listen and serve at port 8081")
	err = http.ListenAndServe(":8081", router)
	if err != nil {
		panic(err)
	}
}
