package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Bachelor-project-f20/eventToGo"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Bachelor-project-f20/api/pkg/graphql"
	sse "github.com/Bachelor-project-f20/api/pkg/sse"
	"github.com/Bachelor-project-f20/shared/config"
	"github.com/go-chi/chi"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
			UseEmitter:  true,
			UseListener: true,
		},
	)
	if err != nil {
		log.Fatalln("configuration failed, error: ", err)
		panic("configuration failed")
	}

	resolver := graphql.Resolver{
		Emitter: configRes.EventEmitter,
	}

	eventChan, err := setupEventListener(configRes.EventListener)
	if err != nil {
		log.Fatalln("setup eventlistener failed, error: ", err)
		panic("setup eventlistener failed")
	}

	router := chi.NewRouter()

	// Add CORS middleware around every request
	// See https://github.com/rs/cors for full option listing
	//router.Use(cors.AllowAll().Handler)
	router.Use(cors.New(cors.Options{
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"X-Requested-With", "Accept", "Authorization", "Accept-Language", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}).Handler)

	//srv := handler.NewDefaultServer(graphql.NewExecutableSchema(graphql.Config{Resolvers: &resolver}))
	srv := handler.New(graphql.NewExecutableSchema(graphql.Config{Resolvers: &resolver}))

	//srv.AddTransport(sse.SSE{})
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})

	srv.SetQueryCache(lru.New(1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New(100),
	})

	router.Handle("/", playground.Handler("GraphQL Playground", "/api"))
	router.Handle("/api", srv)

	sseHandler := sse.NewSSEHandler(eventChan)
	router.HandleFunc("/sse", sseHandler.Handler)

	go func() {
		fmt.Println("Serving metrics API")

		h := http.NewServeMux()
		h.Handle("/metrics", promhttp.Handler())

		http.ListenAndServe(":9191", h)
	}()

	log.Println("API: Listen and serve at port 8081")
	err = http.ListenAndServe(":8081", router)
	if err != nil {
		panic(err)
	}
}
