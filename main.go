package main

import (
	"log"

	"github.com/ilhamsyahids/owldetect/server"
)

func main() {
	// init server
	server.InitHandler()

	// define port, we need to set it as env for Heroku deployment
	port := server.GetPort()

	// run server
	err := server.ServeServer(port)
	if err != nil {
		log.Fatalf("Unable to run server due: %v", err)
	}
}
