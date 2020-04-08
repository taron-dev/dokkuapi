package main

import (
	log "github.com/ondro2208/dokkuapi/logger"
	"github.com/ondro2208/dokkuapi/server"
	str "github.com/ondro2208/dokkuapi/store"
)

func main() {
	store, err := str.NewStore()
	if err != nil {
		log.ErrorLogger.Fatal("Can't create store")
	}

	s := server.NewServer(store)
	s.Start()
}
