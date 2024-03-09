package main

import (
	"context"
	"log"
	"os"

	"github.com/tedmo/go-rest-template/internal/http"
	"github.com/tedmo/go-rest-template/internal/postgres"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
	os.Exit(0)
}

func run() error {
	ctx := context.Background()

	db, err := postgres.NewDBFromEnv(ctx)
	if err != nil {
		return err
	}
	defer db.Close()

	server := &http.Server{
		Port:        8080,
		UserService: postgres.NewUserRepo(db),
	}

	return server.ListenAndServe()
}
