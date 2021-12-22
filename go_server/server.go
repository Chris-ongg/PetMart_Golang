package main

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go_server/petMart_server/database"
	"go_server/petMart_server/dependencies/implementation"
	"go_server/petMart_server/generated"
	"go_server/petMart_server/middleware"
	"go_server/petMart_server/resolvers"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

const defaultPort = "8080"

func main() {
	err := godotenv.Load(".env")

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}

	database.InitialiseDB(client)

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{

		AllowedOrigins:   []string{"http://localhost:3000"}, // Use this to allow specific origin hosts
		//AllowedOrigins:   []string{"https://*", "http://*"},
		AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	router.Use(middleware.InjectHttpMiddleware())

	QueryServices := implementation.NewQueryService(client)
	MutationServices := implementation.NewMutationService(client)

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolvers.Resolver{
		QueryService: QueryServices,
		MutationService: MutationServices,
	}}))

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/graphql", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port,router))
}
