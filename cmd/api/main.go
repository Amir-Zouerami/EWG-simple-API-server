package main

import (
	"log"

	"github.com/Amir-Zouerami/EWG-simple-API-server/internal/env"
	"github.com/Amir-Zouerami/EWG-simple-API-server/internal/store"
)

func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
	}

	store := store.NewStorage(nil)

	app := &application{
		config: cfg,
		store:  store,
	}

	mux := app.mount()

	log.Fatal(app.run(mux))
}
