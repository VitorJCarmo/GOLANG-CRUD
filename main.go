package main

import (
	"log"
)

func main() {
	store, err := NewPostgresStore()
	if err != nil {
		log.Fatal(err.Error())
	}

	if err := store.Init(); err != nil {
		log.Fatal(err.Error())
	}
	server := NewAPIServer(":8080", store)
	server.Run()
}
