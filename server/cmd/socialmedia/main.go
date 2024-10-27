package main

import (
	"log"

	"github.com/maliByatzes/socialmedia/config"
	"github.com/maliByatzes/socialmedia/postgres"
	"github.com/maliByatzes/socialmedia/http"
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

  srv := http.NewServer(db)
  defer srv.Close()
  log.Fatal(srv.Run(cfg.Port))
}
