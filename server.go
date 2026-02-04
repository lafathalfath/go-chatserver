package main

import (
	"context"
	"errors"
	contextkeys "github.com/lafathalfath/go-chatserver/context-keys"
	"github.com/lafathalfath/go-chatserver/database"
	"github.com/lafathalfath/go-chatserver/graph"
	"github.com/lafathalfath/go-chatserver/graph/resolvers"
	"github.com/lafathalfath/go-chatserver/helpers"
	"github.com/lafathalfath/go-chatserver/middlewares"
	"log"
	"net/http"
	"os"
	"slices"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
	"github.com/vektah/gqlparser/v2/ast"
)

const defaultPort = "8080"

func protocol() string {
	isHttps := helpers.Env("ENV") != "development"
	if isHttps {
		return "https://"
	}
	return "http://"
}

func main() {
	database.Connect()

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	serverAddress := protocol() + helpers.Env("HOST") + ":" + helpers.Env("PORT")

	router := chi.NewRouter()
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{helpers.Env("CLIENT_ADDRESS"), serverAddress},
		AllowCredentials: true,
		Debug:            true,
	}).Handler)

	srv := handler.New(graph.NewExecutableSchema(graph.Config{
		Resolvers: &resolvers.Resolver{},
		Directives: graph.DirectiveRoot{
			Auth: graph.AuthDirective,
		},
	}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				origin := r.Header.Get("Origin")
				if origin == "" || origin == r.Header.Get("Host") {
					return true
				}
				return slices.Contains([]string{helpers.Env("CLIENT_ADDRESS")}, origin)
			},
		},
		InitFunc: func(
			ctx context.Context,
			initPayload transport.InitPayload,
		) (context.Context, *transport.InitPayload, error) {

			req, ok := ctx.Value(contextkeys.ContextKeyRequest).(*http.Request)
			if !ok || req == nil {
				return ctx, &initPayload, errors.New("no http request")
			}

			cookie, err := req.Cookie("accessToken")
			if err != nil {
				return ctx, &initPayload, errors.New("unauthorized")
			}

			userId, err := helpers.ParseAccessToken(cookie.Value)
			if err != nil {
				return ctx, &initPayload, errors.New("invalid token")
			}

			ctx = context.WithValue(ctx, contextkeys.UserIDContextKey, userId)
			return ctx, &initPayload, nil
		},
	})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query",
		middlewares.ContextMiddleware(
			middlewares.AuthMiddleware(srv),
		),
	)

	log.Printf("connect to %s/ for GraphQL playground", serverAddress)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
