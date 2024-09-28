package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/shivajichalise/crud/internal/database"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	dbURL := "postgres://<username>:<password>@localhost:5432/coursesDB?sslmode=disable"

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("ERROR: Can't connect to DB")
	}

	queries := database.New(conn)
	apiConf := apiConfig{
		DB: queries,
	}

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Post("/persons", apiConf.handlerCreatePerson)
	v1Router.Put("/persons/{name}", apiConf.handlerUpdatePerson)
	v1Router.Get("/persons", apiConf.handlerGetPersons)
	v1Router.Get("/persons/{name}", apiConf.handlerGetPersonByName)
	v1Router.Delete("/persons/{name}", apiConf.handlerDeletePerson)

	router.Mount("/v1", v1Router)

	server := &http.Server{
		Handler: router,
		Addr:    ":8000",
	}

	fmt.Println("INFO: Server is starting on port: ", 8000)
	server_err := server.ListenAndServe()
	if server_err != nil {
		log.Fatal(server_err)
	}
}
