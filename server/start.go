package server

import "log"

func Start(addr *string) {

	instance := NewInstance(addr)

	log.Println("Server listening at http://localhost:8080")
	instance.Start()

}
