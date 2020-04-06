package api

import (
	"log"
	"net/http"

	sse "github.com/Bachelor-project-f20/api/sse"
	"github.com/Bachelor-project-f20/eventToGo"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Bachelor-project-f20/api/graphql"
	"github.com/Bachelor-project-f20/shared/config"
	"github.com/go-chi/chi"
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

	srv.AddTransport(sse.SSE{})
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})

	srv.SetQueryCache(lru.New(1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New(100),
	})

	//return srv
	// srv.AddTransport(transport.Websocket{
	// 	Upgrader: websocket.Upgrader{
	// 		CheckOrigin: func(r *http.Request) bool {
	// 			// Check against your desired domains here
	// 			fmt.Println("HERE")
	// 			return r.Host == "http://localhost:8081"
	// 		},
	// 		ReadBufferSize:  1024,
	// 		WriteBufferSize: 1024,
	// 	},
	// })

	//
	// srv.AroundResponses(func(ctx context.Context, next gql.ResponseHandler) *gql.Response {
	// 	// This function will be called around each response in the operation. next() will evaluate
	// 	// and return a single response.
	// 	fmt.Println("HERE!!!!!!")
	// 	s := next(ctx)
	// 	fmt.Println(s)
	// 	return s
	// })

	// srv.AroundOperations(func(ctx context.Context, next gql.OperationHandler) gql.ResponseHandler {
	// 	// This function will be called around each response in the operation. next() will evaluate
	// 	// and return a single response.
	// 	fmt.Println("HERE!!!!!!")
	// 	s := next(ctx)
	// 	fmt.Println(s)
	// 	return s
	// })

	router.Handle("/", playground.Handler("GraphQL Playground", "/api"))
	router.Handle("/api", srv)

	log.Println("API: Listen and serve at port 8081")
	err = http.ListenAndServe(":8081", router)
	if err != nil {
		panic(err)
	}
}
