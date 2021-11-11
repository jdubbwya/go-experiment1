package server

import (
	"log"
)

func Start(addr *string) Instance {

	instance := NewInstance(addr)


	defer instance.Start()
	log.Println("Server listening at http://localhost:8080")

	return instance
}
