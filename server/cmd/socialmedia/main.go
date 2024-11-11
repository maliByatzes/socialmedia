package main

import (
	"log"

	"github.com/maliByatzes/socialmedia/config"
	"github.com/maliByatzes/socialmedia/http"
	"github.com/maliByatzes/socialmedia/postgres"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("cannot create new config: %v", err)
	}

	db := postgres.NewDB(cfg.DBURL)
	if err := db.Open(); err != nil {
		log.Fatalf("cannot open database: %v", err)
	}
	defer db.Close()

	srv, err := http.NewServer(db, cfg.SecretKey)
	if err != nil {
		log.Fatal(err)
	}
	defer srv.Close()
	log.Fatal(srv.Run(cfg.Port))
}
