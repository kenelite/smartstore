package main

import (
	"log"

	"github.com/kenelite/smartstore/internal/app"
	"github.com/kenelite/smartstore/internal/config"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	srv := app.NewHTTPServer(cfg)
	log.Printf("smartstore gateway listening on %s (env=%s)", cfg.HTTP.Addr, cfg.Env)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
